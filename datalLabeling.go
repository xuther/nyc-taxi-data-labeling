package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var numberStarted int
var numberSaved int
var numberFailed int
var numberDropped int

var blocks []Block
var tracts []Tract
var counties []County

var Config configuration
var BlocksToTracts map[string][]int
var TractsToCounties map[string][]int

func main() {
	var configAddress = flag.String("configAddress", "./config.json", "Configuration File")

	Config = importConfig(*configAddress)

	var err error
	counties, tracts, blocks, err = getCensusData()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %v couties\nFound %v tracts\nFound %v blocks", len(counties), len(tracts), len(blocks))

	BlocksToTracts = mapBlocksToTracts()
	fmt.Printf("\nMapped %v tracts\n", len(BlocksToTracts))

	TractsToCounties = mapTractsToCounties()
	fmt.Printf("\nMapped %v counties\n", len(TractsToCounties))

	fmt.Printf("Beginning Labelling ")
	fmt.Printf("Using %v max processors.\n", runtime.NumCPU())
	fmt.Printf("Using %v labeling routines.\n", Config.LabelingRoutines)

	channelSize := 250
	//make our channels
	//Channel for unlabeled data to be channeled
	preChannel := make(chan []string, channelSize)
	//Channel for labeled data to be saved
	postChannel := make(chan []string, channelSize)
	//channel to save dropped records for later
	droppedChannel := make(chan []string, 25)
	//channel to save failed records for later
	failedChannel := make(chan []string, 25)

	var wg sync.WaitGroup
	wg.Add(1)

	go readIntoChannel(preChannel, postChannel)
	go finishAndSave(postChannel, wg)
	go saveFailed(failedChannel)
	go saveDropped(droppedChannel)

	for i := 0; i < Config.LabelingRoutines; i++ {
		go labeler(preChannel, postChannel, failedChannel, droppedChannel)
	}

	oldStarted := 0
	oldSaved := 0
	oldDropped := 0
	oldFailed := 0

	for {
		time.Sleep(time.Second * 5)
		log.Printf("%v rows started.", numberStarted)
		log.Printf("%v rows finished.", numberSaved)
		log.Printf("%v rows dropped.", numberDropped)
		log.Printf("%v rows failed.", numberFailed)

		log.Printf("%v rows started in last 5 seconds", numberStarted-oldStarted)
		log.Printf("%v rows finished in last 5 seconds", numberSaved-oldSaved)
		log.Printf("%v rows dropped in last 5 seconds", numberDropped-oldDropped)
		log.Printf("%v rows failed in last 5 seconds", numberFailed-oldFailed)

		//For profiling/debugging
		log.Printf("%v items in the input queue.", len(preChannel))
		log.Printf("%v items in the output queue.\n", len(postChannel))

		oldStarted = numberStarted
		oldSaved = numberSaved
		oldDropped = numberDropped
		oldFailed = numberFailed
		//end if we're all done?
		if numberFailed+numberSaved+numberDropped == numberStarted {
			break
		}
	}
	log.Printf("Done.\n")
	fmt.Printf("Final Report: \n")
	fmt.Printf("%v rows started.\n", numberStarted)
	fmt.Printf("%v rows finished.\n", numberSaved)
	fmt.Printf("%v rows dropped.\n", numberDropped)
	fmt.Printf("%v rows failed.\n", numberFailed)
}

func labeler(inChan <-chan []string, outChan chan<- []string, failedChannel chan<- []string, droppedChannel chan<- []string) {
	for {
		curRow := <-inChan
		//check for null values in columns we care about
		for i := range Config.IndiciesToKeep {
			if len(curRow[Config.IndiciesToKeep[i]]) == 0 {
				numberDropped++
				//add to our list of dropped rows
				droppedChannel <- curRow
				continue
			}
		}
		newRow, err := labelValue(curRow)
		if err != nil {
			curRow = append(curRow, err.Error())
			failedChannel <- curRow
			numberFailed++
			continue
		}

		if len(newRow) != 0 {
			outChan <- newRow
		} else {
			//add to our list of failed rows
			failedChannel <- curRow
			numberFailed++
		}
	}
}

func saveFailed(failedChan <-chan []string) {
	f, err := os.OpenFile(Config.FailedOutputAddress, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error opening failed file:\n")
		log.Printf("Error details %v:\n", err.Error())
		panic(err)
	}

	for {
		toWrite := <-failedChan
		strToWrite := ""
		for i := range toWrite {
			if i == 0 {
				strToWrite += toWrite[i]
			} else {
				strToWrite += "," + toWrite[i]
			}
		}
		_, err := f.WriteString(strToWrite + "\n")
		if err != nil {
			log.Printf("Error writing failed record to file\n")
			log.Printf("Error Details: %v\n", err.Error())
			log.Printf("Record Details: %v\n", toWrite)
		}
	}
}

func saveDropped(droppedChan <-chan []string) {
	f, err := os.OpenFile(Config.FailedOutputAddress, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error opening dropped file:\n")
		log.Printf("Error details %v:\n", err.Error())
		panic(err)
	}

	for {
		toWrite := <-droppedChan
		strToWrite := ""
		for i := range toWrite {
			if i == 0 {
				strToWrite += toWrite[i]
			} else {
				strToWrite += "," + toWrite[i]
			}
		}
		_, err := f.WriteString(strToWrite + "\n")
		if err != nil {
			log.Printf("Error writing dropped record to file\n")
			log.Printf("Error Details: %v\n", err.Error())
			log.Printf("Record Details: %v\n", toWrite)
		}
	}
}

func readIntoChannel(outChan chan<- []string, headerChan chan<- []string) {
	log.Printf("Using input file: %v\n", Config.InputAddress)
	f, err := os.Open(Config.InputAddress)
	if err != nil {
		log.Printf("Error opening input file:\n")
		log.Printf("Error details %v:\n", err.Error())
		panic(err)
	}

	reader := csv.NewReader(f)

	//read the header data in first.
	row, err := reader.Read()
	if err != nil {
		fmt.Printf("Error reading csv headers.")
		panic(err)
	}
	headerChan <- row

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Error reading in csv: %v\n", err.Error())
			fmt.Printf("%v lines read before error \n", numberStarted)
			panic(err)
		}
		if len(row) == 0 {
			continue
		}
		outChan <- row
		numberStarted++
	}

}

//, outChan chan<- []string, inChan <-chan []string

func finishAndSave(labeledRows <-chan []string, wg sync.WaitGroup) {
	defer wg.Done()

	f, err := os.OpenFile(Config.OutputAddress, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error opening output file:\n")
		log.Printf("Error details %v:\n", err.Error())
		panic(err)
	}
	defer f.Close()

	//the first thing through the channel will alwasy be the headers
	oldHeaders := <-labeledRows

	//save headers to the file
	headers := ""
	for i, v := range Config.IndiciesToKeep {
		if i == 0 {
			headers += oldHeaders[v]
			continue
		}
		headers += "," + oldHeaders[v]
	}
	//add the added fields
	for _, v := range Config.AddedFeatures {
		headers += "," + v
	}
	f.WriteString(headers + "\n")

	toWrite := ""
	//loop through and pull rows out of the channel,
	//save out things we care about and continue
	for {
		curRow := <-labeledRows

		toWrite = ""
		for i := range Config.IndiciesToKeep {
			if i == 0 {
				toWrite += curRow[Config.IndiciesToKeep[i]]
			} else {
				toWrite += "," + curRow[Config.IndiciesToKeep[i]]
			}
		}
		for i := range Config.AddedFeatures {
			toWrite += "," + curRow[Config.OriginalFeaturesetSize+i]
		}
		f.WriteString(toWrite + "\n")
		numberSaved++
	}
}

func labelValue(row []string) ([]string, error) {

	//convert to the float64,abort if we can't do it.

	//Do the start coordinates
	x, err := strconv.ParseFloat(row[Config.StartX], 64)
	if err != nil {
		return []string{}, errors.New("Unable to parse StartX")
	}

	y, err := strconv.ParseFloat(row[Config.StartY], 64)
	if err != nil {
		return []string{}, errors.New("Unable to parse StartY")
	}

	county, err := findCounty(x, y)
	if err != nil {
		return []string{}, errors.New("Unable to find Start County")
	}
	row = append(row, county)

	tract, err := findTract(x, y, county)
	if err != nil {
		return []string{}, errors.New("Unable to find Start Tract")
	}
	row = append(row, tract)

	block, err := findBlock(x, y, county, tract)
	if err != nil {
		return []string{}, errors.New("Unable to find Start Block")
	}
	row = append(row, block)

	//Find end information
	x, err = strconv.ParseFloat(row[Config.EndX], 64)
	if err != nil {
		return []string{}, errors.New("Unable to parse EndX")
	}

	y, err = strconv.ParseFloat(row[Config.EndY], 64)
	if err != nil {
		return []string{}, errors.New("Unable to parse EndY")
	}
	county, err = findCounty(x, y)
	if err != nil {
		return []string{}, errors.New("Unable to find End County")
	}
	row = append(row, county)

	tract, err = findTract(x, y, county)
	if err != nil {
		return []string{}, errors.New("Unable to find End Tract")
	}
	row = append(row, tract)

	block, err = findBlock(x, y, county, tract)
	if err != nil {
		return []string{}, errors.New("Unable to find End Block")
	}
	row = append(row, block)

	return row, nil

}

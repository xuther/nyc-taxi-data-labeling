package main

import (
	"encoding/json"
	"io/ioutil"
)

func getCensusData() (counties []County, tracts []Tract, blocks []Block, err error) {
	baseLocation := Config.JsonFileLocation + "/"
	b, err := ioutil.ReadFile(baseLocation + "countyData.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &counties)
	if err != nil {
		return
	}

	b, err = ioutil.ReadFile(baseLocation + "tractData.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &tracts)
	if err != nil {
		return
	}

	b, err = ioutil.ReadFile(baseLocation + "blockData.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &blocks)
	if err != nil {
		return
	}

	return
}

func importConfig(path string) (config configuration) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &config)

	return
}

//returns a map of countyID's to the index of tracts in the
func mapBlocksToTracts() (tractToBlock map[string][]int) {
	tractToBlock = make(map[string][]int)

	for i := range blocks {
		if _, ok := tractToBlock[blocks[i].CountyID+"-"+blocks[i].TractID]; ok {
			tractToBlock[blocks[i].CountyID+"-"+blocks[i].TractID] = append(tractToBlock[blocks[i].CountyID+"-"+blocks[i].TractID], i)
		} else {
			tractToBlock[blocks[i].CountyID+"-"+blocks[i].TractID] = []int{i}
		}
	}
	return
}

//returns a map of countyID's to the index of tracts in the
func mapTractsToCounties() (countyToTract map[string][]int) {
	countyToTract = make(map[string][]int)
	for i := range tracts {

		if _, ok := countyToTract[tracts[i].CountyID]; ok {
			countyToTract[tracts[i].CountyID] = append(countyToTract[tracts[i].CountyID], i)
		} else {
			countyToTract[tracts[i].CountyID] = []int{i}
		}
	}
	return
}

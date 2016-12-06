package main

import (
	"errors"
	"math"
)

func findCounty(x float64, y float64) (string, error) {
	for i := range counties {
		//fmt.Printf("Checking county %v\n", counties[i].CountyID)
		if pointInPoly(x, y, counties[i].Points) {
			return counties[i].CountyID, nil
		}
	}

	return "", errors.New("Unable to find count\n")
}

func findTract(x float64, y float64, county string) (string, error) {
	//log.Printf("Placing point in county %v\n", county)

	for _, i := range TractsToCounties[county] {
		//log.Printf("Checking tract %v.\n", tracts[i].TractID)

		if pointInPoly(x, y, tracts[i].Points) {
			return tracts[i].TractID, nil
		}
	}

	return "", errors.New("Unable to find Tract\n")
}

func findBlock(x float64, y float64, county string, tract string) (string, error) {
	for _, i := range BlocksToTracts[county+"-"+tract] {
		//log.Printf("Checking block %v.\n", blocks[i].BlockID)
		if pointInPoly(x, y, blocks[i].Points) {
			return blocks[i].BlockID, nil
		}
	}

	return "", errors.New("Unable to find block\n")
}

//Is there a better way?
func pointInPoly(x float64, y float64, polygon [][]float64) bool {

	p1x := polygon[0][0]
	p1y := polygon[0][1]

	n := len(polygon)
	inside := false

	var xints float64
	count := 0
	for i := 0; i < n+1; i++ {

		p2x := polygon[i%n][0]
		p2y := polygon[i%n][1]

		if y > math.Min(p1y, p2y) {
			if y <= math.Max(p1y, p2y) {
				if x <= math.Max(p1x, p2x) {
					if p1y != p2y {
						xints = (y-p1y)*(p2x-p1x)/(p2y-p1y) + p1x
					}
					if p1x == p2x || x <= xints {
						count = count + 1
						inside = !inside
						//fmt.Printf("Crossed Line (%v, %v), (%v,%v)\n", p1x, p1y, p2x, p2y)
					}
				}
			}
		}
		p1x = p2x
		p1y = p2y
	}
	if count > 0 {
		//fmt.Printf("Crossed %v lines\n", count)
	}

	//return count%2 != 0
	return inside
}

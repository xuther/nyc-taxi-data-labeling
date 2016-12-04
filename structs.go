package main

//County comes from the 2010 Census Shapefile for counties in NYC
type County struct {
	StateID                  string      `json:"STATEFP10"`
	CountyID                 string      `json:"COUNTYFP10"`
	CountyANSICode           string      `json:"COUNTYNS10"`
	LongID                   string      `json:"GEOID10"`
	Name                     string      `json:"NAME10"`
	ShortDesc                string      `json:"NAMELSAD10"`
	StatsDesc                string      `json:"LSAD10"`
	CountyClass              string      `json:"CLASSFP10"`
	TigerClassCode           string      `json:"MTFCC10"`
	StatisticalAreaCode      string      `json:"CSAFP10"`
	MetropolitanAreaCode     string      `json:"CBSAFP10"`
	MetropolitanDivisionCode string      `json:"METDIVFP10"`
	FunctionalStatus         string      `json:"FUNCSTAT10"`
	LandArea                 int         `json:"ALAND10"`
	WaterArea                int         `json:"AWATER10"`
	InternalLatitudePoint    string      `json:"INTPTLAT10"`
	InternalLongitudePoint   string      `json:"INTPTLON10"`
	Points                   [][]float64 `json:"points"`
}

//Tract comes from the 2010 Census Shapefil for CensusTracts in NYC
type Tract struct {
	StateID                string      `json:"STATEFP10"`
	CountyID               string      `json:"COUNTYFP10"`
	TractID                string      `json:"TRACTCE10"`
	LongID                 string      `json:"GEOID10"`
	Name                   string      `json:"NAME10"`
	ShortDesc              string      `json:"NAMELSAD10"`
	TigerClassCode         string      `json:"MTFCC10"`
	FunctionalStatus       string      `json:"FUNCSTAT10"`
	LandArea               int         `json:"ALAND10"`
	WaterArea              int         `json:"AWATER10"`
	InternalLatitudePoint  string      `json:"INTPTLAT10"`
	InternalLongitudePoint string      `json:"INTPTLON10"`
	Points                 [][]float64 `json:"points"`
}

//Block comes from the 2010 Census Shapefil for CensusBlocks in NYC
type Block struct {
	StateID                string      `json:"STATEFP10"`
	CountyID               string      `json:"COUNTYFP10"`
	TractID                string      `json:"TRACTCE10"`
	BlockID                string      `json:"BLOCKCE10"`
	LongID                 string      `json:"GEOID10"`
	Name                   string      `json:"NAME10"`
	TigerClassCode         string      `json:"MTFCC10"`
	UrbanRural             string      `json:"UR10"`
	UrbanAreaCode          string      `json:"UACE10"`
	UrbanAreaType          string      `json:"UATYP10"`
	FunctionalStatus       string      `json:"FUNCSTAT10"`
	LandArea               int         `json:"ALAND10"`
	WaterArea              int         `json:"AWATER10"`
	InternalLatitudePoint  string      `json:"INTPTLAT10"`
	InternalLongitudePoint string      `json:"INTPTLON10"`
	Points                 [][]float64 `json:"points"`
}

type configuration struct {
	IndiciesToKeep         []int    `json:"indicies-to-keep"`
	StartX                 int      `json:"start-x"`
	StartY                 int      `json:"start-y"`
	EndX                   int      `json:"end-x"`
	EndY                   int      `json:"end-y"`
	InputAddress           string   `json:"input-file"`
	OutputAddress          string   `json:"output-file"`
	AddedFeatures          []string `json:"features-added"`
	OriginalFeaturesetSize int      `json:"original-featureset-size"`
	LabelingRoutines       int      `json:"labeling-routines"`
	FailedOutputAddress    string   `json:"failed-output-file"`
	DroppedOutputAddress   string   `json:"dropped-output-file"`
}

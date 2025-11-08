package types

type Location struct {
	Latitude string `json:"lat"`
	Longitude string `json:"log"`
	Speed *string `json:"speed"`
}

type PlantGPSDataField struct {
	GpsUTC string
	Lat string
	Long string
	Alt int
	Spd float32
	SpdAcc int
	Head int
	PDOP int
	PosAcc int
	GpsStat int
	FType int
}

type PlantOdoField struct {
	Odo int
	RH int
	FType int
}

type PlantDeviceField struct {
	DIn int
	DOut int
	DevStat int
	FType int
}

type PlantAnalogueDataField struct {
	AnalogueData map[int]int
	FType int
}

type TrackingPlantDataRecord struct {
	SeqNo int
	Reason int
	DateUTC string
	Fields []any
}

type TrackingPlantData struct {
	SerNo string
	IMEI string
	ProdId int
	FW string
	Records []TrackingPlantDataRecord
}

package internal

import (
	"script-tracking-gis/pkg/utils"
	types "script-tracking-gis/types"
	"time"
)

type BuildTrackingPlantDataBuilder struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"log"`
	Speed     *string `json:"speed"`
	SeqNo int
	SerNo string
}

// func BuildTrackingPlantData(scenario types.Scenario) types.TrackingPlantData {
func BuildTrackingPlantData(params BuildTrackingPlantDataBuilder) types.TrackingPlantData {
	currentTime := time.Now().UTC().Format(time.RFC3339Nano)
	trackingPlantData := types.TrackingPlantData{
		SerNo: params.SerNo,
		IMEI: params.SerNo,
		ProdId: 1,
		FW: "1.1.1.1",
		Records: []types.TrackingPlantDataRecord{
			{
				SeqNo: params.SeqNo,
				Reason: 1,
				DateUTC: currentTime,
				Fields: []any{
					types.PlantGPSDataField{
						GpsUTC: currentTime,
						Lat: params.Latitude,
						Long: params.Longitude,
						Alt: 18,
						Spd: utils.ConvertKMHToMS(50),
						SpdAcc: 2,
						Head: 159,
						PDOP: 28,
						PosAcc: 6,
						GpsStat: 3,
						FType: 0,
					},
					types.PlantOdoField{
						Odo: 111,
						RH: 12345,
						FType: 1,
					},
					types.PlantDeviceField{
						DIn: 0,
						DOut: 0,
						DevStat: 2,
						FType: 2,
					},
					types.PlantAnalogueDataField{
						AnalogueData: map[int]int{
							4: utils.RandomNumberFromRange(15, 28),
							1: 4144,
							2: 30,
							5: 3,
							3: 1467,
						},
						FType: 1,
					},
				},
			},
		},
	}

	return trackingPlantData
}

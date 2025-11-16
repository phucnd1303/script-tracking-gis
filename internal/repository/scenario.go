package repository

import (
	"fmt"

	types "github.com/phucnd1303/script-tracking-gis/types"
)

const (
	stageLocal = "local"
	stageDEV   = "dev"
	stageUAT   = "uat"
	stageSIT   = "sit"
)

const (
	sydneyLocationsPath         = "templates/australia/sydney.json"
	duongDinhNgheLocationsPath  = "templates/vietnam/duong-dinh-nghe.danang.json"
	nhaTrangStreetLocationsPath = "templates/vietnam/tran-phu-nha-trang-street.json"
)

func GetScenarios() ([]types.Scenario, error) {
	sydneyLocations, err := LoadLocationsTemplate(sydneyLocationsPath)

	if err != nil {
		fmt.Printf("Error loading locations template: %v", err)

		return nil, err
	}

	duongDinhNgheLocations, err := LoadLocationsTemplate(duongDinhNgheLocationsPath)

	if err != nil {
		fmt.Printf("Error loading locations template: %v", err)

		return nil, err
	}

	scenarios := []types.Scenario{
		{
			Env:            stageDEV,
			SerNo:          "1104222",
			Locations:      duongDinhNgheLocations,
			DelayTime:      1,
			TotalTime:      100000,
			DelayStartTime: 0,
		},
		{
			Env:            stageDEV,
			SerNo:          "1107056",
			Locations:      duongDinhNgheLocations,
			DelayTime:      1,
			TotalTime:      100000,
			DelayStartTime: 5,
		},
		{
			Env:            stageDEV,
			SerNo:          "1180901",
			Locations:      duongDinhNgheLocations,
			DelayTime:      1,
			TotalTime:      100000,
			DelayStartTime: 10,
		},
		{
			Env:            stageDEV,
			SerNo:          "12398",
			Locations:      sydneyLocations,
			DelayTime:      1,
			TotalTime:      100000,
			DelayStartTime: 0,
		},
		{
			Env:            stageDEV,
			SerNo:          "8102025",
			Locations:      sydneyLocations,
			DelayTime:      1,
			TotalTime:      100000,
			DelayStartTime: 5,
		},
	}

	return scenarios, nil
}

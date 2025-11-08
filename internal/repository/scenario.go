package repository

import (
	"fmt"
	types "script-tracking-gis/types"
)

const (
	stageLocal = "local"
	stageDEV   = "dev"
	stageUAT   = "uat"
	stageSIT   = "sit"
)

const (
	sydneyLocationsPath        = "templates/australia/sydney.json"
	duongDinhNgheLocationsPath = "templates/vietnam/duong-dinh-nghe.danang.json"
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
			Env:            stageLocal,
			SerNo:          "1234567890",
			Locations:      sydneyLocations,
			DelayTime:      3,
			TotalTime:      200,
			DelayStartTime: 0,
		},
		{
			Env:            stageLocal,
			SerNo:          "111111",
			Locations:      duongDinhNgheLocations,
			DelayTime:      1,
			TotalTime:      200,
			DelayStartTime: 5,
		},
	}

	return scenarios, nil
}

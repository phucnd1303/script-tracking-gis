package repository

import (
	"fmt"

	"github.com/phucnd1303/script-tracking-gis/internal/models"
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

func GetScenarios() ([]models.Scenario, error) {
	nhaTrangStreetLocations, err := LoadLocationsTemplate(nhaTrangStreetLocationsPath)

	if err != nil {
		fmt.Printf("Error loading locations template: %v", err)

		return nil, err
	}

	scenarios := []models.Scenario{
		{
			Env:            stageSIT,
			SerNo:          "1104222",
			Locations:      nhaTrangStreetLocations,
			DelayTime:      1,
			TotalTime:      100000,
			DelayStartTime: 0,
		},
		{
			Env:            stageSIT,
			SerNo:          "5112025",
			Locations:      nhaTrangStreetLocations,
			DelayTime:      1,
			TotalTime:      100000,
			DelayStartTime: 5,
		},
		{
			Env:            stageSIT,
			SerNo:          "8102025",
			Locations:      nhaTrangStreetLocations,
			DelayTime:      1,
			TotalTime:      100000,
			DelayStartTime: 10,
		},
	}

	return scenarios, nil
}

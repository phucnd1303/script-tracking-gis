package internal

import (
	types "script-tracking-gis/types"
	utils "script-tracking-gis/utils"
	"time"
)

func reverseLocations(locations *[]types.Location) {
	for firstIndex, lastIndex := 0, len(*locations) - 1; firstIndex < lastIndex; firstIndex, lastIndex = firstIndex + 1, lastIndex -1 {
		(*locations)[firstIndex], (*locations)[lastIndex] = (*locations)[lastIndex], (*locations)[firstIndex]
	}
}

func PlantTracking(scenario types.Scenario) {
	apiKey := utils.GetAPIKey(scenario.Env)
	apiURL := utils.GetPlantAPIURL(scenario.Env)
	counter := 0
	locationIndex := 0
	isReverse := false
	// why I can't input directly number?
	ticker := time.NewTicker(time.Duration(scenario.DelayTime) * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		counter++

		if counter >= scenario.TotalTime {
			return
		}

		// locations := scenario.Locations
		locations := make([]types.Location, len(scenario.Locations))
		copy(locations, scenario.Locations)

		if isReverse {
			reverseLocations(&locations)
		}

		SendTrackingPlantData(SendTrackingPlantDataBuilder{
			State: scenario.Env,
			Location: locations[locationIndex],
			SeqNo: counter,
			SerNo: scenario.SerNo,
			APIKey: apiKey,
			APIURL: apiURL,
		})

		// fmt.Printf("[locationIndex %d - SerNo %s]\n", locationIndex, scenario.SerNo)

		if locationIndex >= len(scenario.Locations) - 1 {
			isReverse = true
			locationIndex = 0
		} else {
			locationIndex++
		}
	}
}

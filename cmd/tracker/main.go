package main

import (
	"fmt"
	"script-tracking-gis/internal/api"
	"script-tracking-gis/internal/config"
	"script-tracking-gis/internal/repository"
	"script-tracking-gis/internal/tracker"
	"script-tracking-gis/types"
	"sync"
	"time"
)

func main() {
	cfg, error := config.NewConfig()

	if error != nil {
		fmt.Printf("Error loading .env file: %v", error)
		return
	}

	scenarios, error := repository.GetScenarios()

	if error != nil {
		fmt.Printf("Error getting scenarios: %v", error)
		return
	}

	apiClient := api.NewClient(cfg)
	tracker := tracker.NewTracker(cfg, apiClient)

	var waitGroup sync.WaitGroup

	for _, scenario := range scenarios {
		waitGroup.Add(1)

		go func(scenario types.Scenario) {
			defer waitGroup.Done()

			time.Sleep(time.Duration(scenario.DelayStartTime) * time.Second)

			// internal.PlantTracking(cfg, scenario)
			tracker.Track(scenario)
		}(scenario)
	}

	waitGroup.Wait()

	fmt.Println("All scenarios completed")
}

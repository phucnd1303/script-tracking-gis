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
	cfg, err := config.NewConfig()

	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		return
	}

	scenarios, err := repository.GetScenarios()

	if err != nil {
		fmt.Printf("Error getting scenarios: %v", err)
		return
	}

	apiClient := api.NewClient(cfg)
	tracker := tracker.NewTracker(cfg, apiClient)

	var wg sync.WaitGroup

	for _, scenario := range scenarios {
		wg.Add(1)

		go func(scenario types.Scenario) {
			defer wg.Done()

			time.Sleep(time.Duration(scenario.DelayStartTime) * time.Second)

			tracker.Track(scenario)
		}(scenario)
	}

	wg.Wait()

	fmt.Println("All scenarios completed")
}

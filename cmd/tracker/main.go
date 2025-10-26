package main

import (
	"fmt"
	internal "script-tracking-gis/internal"
	"script-tracking-gis/types"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	var error = godotenv.Load()

	if error != nil {
		fmt.Printf("Error loading .env file: %v", error)
		return
	}

	scenarios, error := internal.GetScenarios()

	if error != nil {
		fmt.Printf("Error getting scenarios: %v", error)
		return
	}

	var waitGroup sync.WaitGroup

	for _, scenario := range scenarios {
		waitGroup.Add(1)

		go func(scenario types.Scenario) {
			time.Sleep(time.Duration(scenario.DelayStartTime) * time.Second)

			defer waitGroup.Done()

			internal.PlantTracking(scenario)
		}(scenario)
	}

	waitGroup.Wait()

	fmt.Println("All scenarios completed")
}

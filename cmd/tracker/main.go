package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"script-tracking-gis/internal/api"
	"script-tracking-gis/internal/config"
	"script-tracking-gis/internal/repository"
	"script-tracking-gis/internal/tracker"
	"script-tracking-gis/types"
	"sync"
	"syscall"
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

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		fmt.Printf("\n🛑 Received signal: %v\n", sig)
		fmt.Println("⏳ Initiating graceful shutdown...")
		cancel()
	}()

	var wg sync.WaitGroup

	for _, scenario := range scenarios {
		wg.Add(1)

		go func(scenario types.Scenario) {
			defer wg.Done()

			select {
			case <-time.After(time.Duration(scenario.DelayStartTime) * time.Second):
			case <-ctx.Done():
				fmt.Printf("❌ Tracker %s cancelled before start\n", scenario.SerNo)
				return
			}

			fmt.Printf("▶️  Starting tracker for SerNo: %s\n", scenario.SerNo)

			if err := tracker.Track(ctx, scenario); err != nil {
				if err == context.Canceled {
					fmt.Printf("⏹️  Tracker %s stopped gracefully\n", scenario.SerNo)
				} else {
					fmt.Printf("❌ Tracker %s error: %v\n", scenario.SerNo, err)
				}
			}
		}(scenario)
	}

	wg.Wait()

	fmt.Println("All scenarios completed")
}

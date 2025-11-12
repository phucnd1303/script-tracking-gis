package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"script-tracking-gis/internal/api"
	"script-tracking-gis/internal/config"
	"script-tracking-gis/internal/repository"
	"script-tracking-gis/internal/tracker"
	"script-tracking-gis/types"
	"strings"
	"sync"
	"syscall"
	"time"
)

type TrackerController struct {
	cancels map[string]context.CancelFunc
	mu sync.RWMutex
}

func NewTrackerController() *TrackerController {
	return &TrackerController{
		cancels: make(map[string]context.CancelFunc),
	}
}

func (tc *TrackerController) StopTracker(serNo string) bool {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if cancel, ok := tc.cancels[serNo]; ok {
		cancel()
	  delete(tc.cancels, serNo)

		return true
	}

	fmt.Printf("Tracker for SerNo: %s not found\n", serNo)
	return false
}

func (tc *TrackerController) StopAllTrackers() {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	for _, cancel := range tc.cancels {
		cancel()
	}

	tc.cancels = make(map[string]context.CancelFunc)
}

func (tc *TrackerController) AddTracker(serNo string, cancel context.CancelFunc) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.cancels[serNo] = cancel
}

func (tc *TrackerController) DeleteTracker(serNo string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	delete(tc.cancels, serNo)
}

func (tc *TrackerController) ListTrackers() {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	if len(tc.cancels) == 0 {
		fmt.Println("📟 No active trackers")
	} else {
		for serNo := range tc.cancels {
			fmt.Printf("📟 Active Tracker: %s\n", serNo)
		}
	}
}

func startInteractiveMode(controller *TrackerController) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n📟 Interactive Mode Commands:")
	fmt.Println("  list              - List active trackers")
	fmt.Println("  stop <SerNo>      - Stop specific tracker")
	fmt.Println("  stopall           - Stop all trackers")
	fmt.Println("  quit              - Stop all and exit")
	fmt.Println()

	go func() {
		for scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(input)

			if len(parts) == 0 {
				continue
			}

			switch parts[0] {
			case "list":
				controller.ListTrackers()

			case "stop":
				if len(parts) < 2 {
					fmt.Println("Usage: stop <serNo>")
					continue
				}

				serNo := parts[1]
				if controller.StopTracker(serNo) {
					fmt.Printf("Tracker for SerNo: %s stopped\n", serNo)
				} else {
					fmt.Printf("Tracker for SerNo: %s not found\n", serNo)
				}

			case "stopall":
				controller.StopAllTrackers()
				fmt.Println("All trackers stopped")

			case "quit", "exit":
				fmt.Println("Exiting...")
				controller.StopAllTrackers()
				os.Exit(0)

			case "help", "?":
				fmt.Println("\n📟 Available Commands:")
				fmt.Println("  list              - List active trackers")
				fmt.Println("  stop <SerNo>      - Stop specific tracker")
				fmt.Println("  stopall           - Stop all trackers")
				fmt.Println("  quit              - Stop all and exit")

			default:
				fmt.Printf("Unknown command: %s (type 'help' for commands)\n", parts[0])
			}

			fmt.Print("> ")
		}
	}()
}
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

	controller := NewTrackerController()

	startInteractiveMode(controller)

	globalCtx, globalCancel := context.WithCancel(context.Background())
	defer globalCancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		fmt.Printf("\n🛑 Received signal: %v\n", sig)
		fmt.Println("⏳ Initiating graceful shutdown...")
		controller.StopAllTrackers()
		globalCancel()
	}()

	var wg sync.WaitGroup

	for _, scenario := range scenarios {
		wg.Add(1)

		go func(scenario types.Scenario) {
			defer wg.Done()

			trackerCtx, trackerCancel := context.WithCancel(globalCtx)
			defer trackerCancel()

			controller.AddTracker(scenario.SerNo, trackerCancel)
			defer controller.DeleteTracker(scenario.SerNo)

			select {
			case <-time.After(time.Duration(scenario.DelayStartTime) * time.Second):
			case <-trackerCtx.Done():
				fmt.Printf("❌ Tracker %s cancelled before start\n", scenario.SerNo)
				return
			}

			if err := tracker.Track(trackerCtx, scenario); err != nil {
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

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/phucnd1303/script-tracking-gis/internal/api"
	"github.com/phucnd1303/script-tracking-gis/internal/cli"
	"github.com/phucnd1303/script-tracking-gis/internal/config"
	"github.com/phucnd1303/script-tracking-gis/internal/controller"
	"github.com/phucnd1303/script-tracking-gis/internal/monitor"
	"github.com/phucnd1303/script-tracking-gis/internal/repository"
	"github.com/phucnd1303/script-tracking-gis/internal/tracker"
	"github.com/phucnd1303/script-tracking-gis/types"
)

type Application struct {
	config       *config.Config
	controller   *controller.TrackerController
	tracker      *tracker.Tracker
	cli          *cli.InteractiveCLI
	sm           *monitor.SystemMonitor
	globalCtx    context.Context
	globalCancel context.CancelFunc
}

func NewApplication() (*Application, error) {
	cfg, err := config.NewConfig()

	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	apiClient := api.NewClient(cfg)
	tracker := tracker.NewTracker(cfg, apiClient)
	tc := controller.NewTrackerController()
	sm := monitor.NewSystemMonitor()
	cli := cli.NewInteractiveCLI(tc, sm)
	globalCtx, globalCancel := context.WithCancel(context.Background())

	return &Application{
		config:       cfg,
		controller:   tc,
		tracker:      tracker,
		cli:          cli,
		sm:           sm,
		globalCtx:    globalCtx,
		globalCancel: globalCancel,
	}, nil
}

func (app *Application) shutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		fmt.Printf("\n🛑 Received signal: %v\n", sig)
		fmt.Println("⏳ Initiating graceful shutdown...")
		app.controller.StopAllTrackers()
		app.globalCancel()
	}()
}

func (app *Application) Run() error {
	app.cli.Start()
	app.shutdown()
	app.sm.StartAutoMonitor(app.globalCtx, 10*time.Second)

	scenarios, err := repository.GetScenarios()

	if err != nil {
		return fmt.Errorf("error getting scenarios: %w", err)
	}

	var wg sync.WaitGroup

	for _, scenario := range scenarios {
		wg.Add(1)

		go func(scenario types.Scenario) {
			defer wg.Done()

			trackerCtx, trackerCancel := context.WithCancel(app.globalCtx)
			defer trackerCancel()

			app.controller.AddTracker(scenario.SerNo, trackerCancel)
			defer app.controller.DeleteTracker(scenario.SerNo)

			select {
			case <-time.After(time.Duration(scenario.DelayStartTime) * time.Second):
			case <-trackerCtx.Done():
				fmt.Printf("❌ Tracker %s cancelled before start\n", scenario.SerNo)
				return
			}

			if err := app.tracker.Track(trackerCtx, scenario); err != nil {
				if err == context.Canceled {
					fmt.Printf("⏹️ Tracker %s stopped gracefully\n", scenario.SerNo)
				} else {
					fmt.Printf("❌ Tracker %s error: %v\n", scenario.SerNo, err)
				}
			}
		}(scenario)
	}

	wg.Wait()
	fmt.Println("All scenarios completed")

	return nil
}

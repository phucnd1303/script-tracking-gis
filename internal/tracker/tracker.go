package tracker

import (
	"context"
	"fmt"
	"script-tracking-gis/internal/api"
	"script-tracking-gis/internal/config"
	"script-tracking-gis/types"
	"time"
)

type Tracker struct {
	config *config.Config
	apiClient *api.Client
}

func NewTracker(cfg *config.Config, apiClient *api.Client) *Tracker {
	return &Tracker{
		config: cfg,
		apiClient: apiClient,
	}
}

func reverseLocations(locations *[]types.Location) {
	for firstIndex, lastIndex := 0, len(*locations) - 1; firstIndex < lastIndex; firstIndex, lastIndex = firstIndex + 1, lastIndex -1 {
		(*locations)[firstIndex], (*locations)[lastIndex] = (*locations)[lastIndex], (*locations)[firstIndex]
	}
}

func (tracker *Tracker) Track(ctx context.Context, scenario types.Scenario) error {
	counter := 0
	locationIndex := 0
	isReverse := false
	ticker := time.NewTicker(time.Duration(scenario.DelayTime) * time.Second)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			counter++

			if counter >= scenario.TotalTime {
				return nil
			}

			locations := make([]types.Location, len(scenario.Locations))
			copy(locations, scenario.Locations)

			if isReverse {
				reverseLocations(&locations)
			}

			_, err := tracker.apiClient.SendTrackingPlantData(ctx, api.TrackingRequest{
				State:    scenario.Env,
				Location: locations[locationIndex],
				SeqNo:    counter,
				SerNo:    scenario.SerNo,
			})

			if err != nil {
				if err == context.Canceled {
					return err
				}

			  return fmt.Errorf("error sending tracking plant data: %w", err)
			} 
			// else {
			// 	fmt.Printf("Response-body: %s\n", res.Body)
			// 	fmt.Printf("Response-status: %d\n", res.StatusCode)
			// }

			if locationIndex >= len(scenario.Locations)-1 {
				isReverse = true
				locationIndex = 0
			} else {
				locationIndex++
			}
		}
	}
}

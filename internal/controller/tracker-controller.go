package controller

import (
	"context"
	"fmt"
	"sync"
)

type TrackerController struct {
	cancels map[string]context.CancelFunc
	mu      sync.RWMutex
}

func NewTrackerController() *TrackerController {
	return &TrackerController{
		cancels: make(map[string]context.CancelFunc),
	}
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

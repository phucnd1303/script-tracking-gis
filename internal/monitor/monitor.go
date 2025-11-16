package monitor

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type SystemMetrics struct {
	StartTime time.Time
	Uptime    time.Duration

	// Memory stats
	MemoryAlloc uint64 // Current allocated memory
	MemoryTotal uint64 // Total allocated memory
	MemorySys   uint64 // Memory from OS
	NumGC       uint32 // GC runs

	// Goroutine stats
	NumGoroutine int
	NumCPU       int
}

type SystemMonitor struct {
	startTime time.Time
}

func NewSystemMonitor() *SystemMonitor {
	return &SystemMonitor{
		startTime: time.Now(),
	}
}

func (sm *SystemMonitor) GetMetrics() SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemMetrics{
		StartTime:    sm.startTime,
		Uptime:       time.Since(sm.startTime),
		MemoryAlloc:  m.Alloc,
		MemoryTotal:  m.TotalAlloc,
		MemorySys:    m.Sys,
		NumGC:        m.NumGC,
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
	}
}

func (sm *SystemMonitor) PrintMetrics(metrics SystemMetrics) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("📊 SYSTEM MONITOR - %s\n", time.Now().Format("15:04:05"))
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("📱 Application:")
	fmt.Printf("  Uptime:          %v\n", metrics.Uptime.Round(time.Second))

	fmt.Println("\n💻 System:")
	fmt.Printf("  CPU Cores:       %d\n", metrics.NumCPU)
	fmt.Printf("  Goroutines:      %d\n", metrics.NumGoroutine)
	fmt.Printf("  Memory Total:     %d MB\n", metrics.MemoryTotal/1024/1024)
	fmt.Printf("  Memory Alloc:    %d MB\n", metrics.MemoryAlloc/1024/1024)
	fmt.Printf("  Memory Sys:      %d MB\n", metrics.MemorySys/1024/1024)
	fmt.Printf("  GC Runs:         %d\n", metrics.NumGC)

	fmt.Println(strings.Repeat("=", 70))
}

func (sm *SystemMonitor) StartAutoMonitor(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metrics := sm.GetMetrics()
				sm.PrintMetrics(metrics)
			}
		}
	}()
}

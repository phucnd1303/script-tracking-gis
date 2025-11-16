package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/phucnd1303/script-tracking-gis/internal/controller"
	"github.com/phucnd1303/script-tracking-gis/internal/monitor"
)

type InteractiveCLI struct {
	controller *controller.TrackerController
	sm         *monitor.SystemMonitor
}

func NewInteractiveCLI(controller *controller.TrackerController, sm *monitor.SystemMonitor) *InteractiveCLI {
	return &InteractiveCLI{
		controller: controller,
		sm:         sm,
	}
}

func (cli *InteractiveCLI) Start() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n📟 Interactive Mode Commands:")
	fmt.Println("  list              - List active trackers")
	fmt.Println("  monitor           - Monitor system metrics")
	fmt.Println("  stop <SerNo>      - Stop specific tracker")
	fmt.Println("  stopall           - Stop all trackers")
	fmt.Println("  quit              - Stop all and exit")
	fmt.Println()
	fmt.Print("> ")

	go func() {
		for scanner.Scan() {
			command := strings.TrimSpace(scanner.Text())
			cli.manageCommand(command)
			fmt.Print("> ")
		}
	}()
}

func (cli *InteractiveCLI) manageCommand(command string) {
	parts := strings.Fields(command)

	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "list":
		cli.controller.ListTrackers()

	case "monitor":
		metrics := cli.sm.GetMetrics()
		cli.sm.PrintMetrics(metrics)

	case "stop":
		if len(parts) < 2 {
			fmt.Println("Usage: stop <serNo>")
			return
		}

		cli.controller.StopTracker(parts[1])

	case "stopall":
		cli.controller.StopAllTrackers()

	case "quit", "exit":
		fmt.Println("Exiting...")
		cli.controller.StopAllTrackers()
		os.Exit(0)

	case "help", "?":
		fmt.Println("\n📟 Available Commands:")
		fmt.Println("  list              - List active trackers")
		fmt.Println("  monitor           - Monitor system metrics")
		fmt.Println("  stop <SerNo>      - Stop specific tracker")
		fmt.Println("  stopall           - Stop all trackers")
		fmt.Println("  quit              - Stop all and exit")

	default:
		fmt.Printf("Unknown command: %s (type 'help' for commands)\n", parts[0])
	}
}

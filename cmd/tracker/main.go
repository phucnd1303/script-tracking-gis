package main

import (
	"fmt"
	"script-tracking-gis/internal/app"
)

func main() {
	application, err := app.NewApplication()

	if err != nil {
		fmt.Printf("Error creating application: %v", err)
		return
	}

	if err := application.Run(); err != nil {
		fmt.Printf("Error running application: %v", err)
		return
	}
}

package repository

import (
	"encoding/json"
	"fmt"
	"os"
	types "script-tracking-gis/types"
)

func LoadLocationsTemplate(path string) ([]types.Location, error) {
	jsonFile, error := os.ReadFile(path)

	if error != nil {
		fmt.Printf("Error opening file: %v", error)
		return nil, error
	}

	var locations []types.Location
	error = json.Unmarshal(jsonFile, &locations)

	if error != nil {
		fmt.Printf("Error unmarshalling file: %v", error)
		return nil, error
	}

	return locations, nil
}

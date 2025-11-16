package repository

import (
	"encoding/json"
	"fmt"
	"os"

	types "github.com/phucnd1303/script-tracking-gis/types"
)

func LoadLocationsTemplate(path string) ([]types.Location, error) {
	jsonFile, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return nil, err
	}

	var locations []types.Location
	err = json.Unmarshal(jsonFile, &locations)

	if err != nil {
		fmt.Printf("Error unmarshalling file: %v", err)
		return nil, err
	}

	return locations, nil
}

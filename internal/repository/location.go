package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/phucnd1303/script-tracking-gis/internal/models"
)

func LoadLocationsTemplate(path string) ([]models.Location, error) {
	jsonFile, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return nil, err
	}

	var locations []models.Location
	err = json.Unmarshal(jsonFile, &locations)

	if err != nil {
		fmt.Printf("Error unmarshalling file: %v", err)
		return nil, err
	}

	return locations, nil
}

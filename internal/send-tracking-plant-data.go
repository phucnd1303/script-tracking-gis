package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	types "script-tracking-gis/types"
	"time"
)

type SendTrackingPlantDataBuilder struct {
	State string
	Location types.Location
	SeqNo int
	SerNo string
	APIKey string
	APIURL string
}

func SendTrackingPlantData(params SendTrackingPlantDataBuilder) {
	trackingPlantData := BuildTrackingPlantData(BuildTrackingPlantDataBuilder{
		Latitude: params.Location.Latitude,
		Longitude: params.Location.Longitude,
		Speed: params.Location.Speed,
		SeqNo: params.SeqNo,
		SerNo: params.SerNo,
	})
	jsonData, error := json.Marshal(trackingPlantData)
	
	if error != nil {
		fmt.Printf("Error marshalling tracking plant data: %v", error)
		return
	}

	request, error := http.NewRequest("POST", params.APIURL, bytes.NewBuffer(jsonData))

	if error != nil {
		fmt.Printf("Error creating request: %v", error)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", params.APIKey)
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/114.0")

	client := &http.Client{
		Timeout: 25 * time.Second,
	}

	response, error := client.Do(request)

	if error != nil {
		fmt.Printf("Error sending tracking plant data: %v\n", error)
		return
	}

	defer response.Body.Close()

	body, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Printf("Error reading response body: %v\n", error)
		return
	}

	// consider another solution to convert this data
	fmt.Printf("Response-body: %s\n", string(body))
	fmt.Printf("Response-status: %d\n", response.StatusCode)
}

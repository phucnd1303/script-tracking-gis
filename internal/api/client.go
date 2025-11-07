package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"script-tracking-gis/internal/config"
	"script-tracking-gis/pkg/utils"
	"script-tracking-gis/types"
	"time"
)

type Client struct {
	config *config.Config
	httpClient *http.Client
}

type TrackingRequest struct {
	State string
	Location types.Location
	SeqNo int
	SerNo string
}

type BuildTrackingPlantDataBuilder struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"log"`
	Speed     *string `json:"speed"`
	SeqNo int
	SerNo string
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

func buildTrackingPlantData(params BuildTrackingPlantDataBuilder) types.TrackingPlantData {
	currentTime := time.Now().UTC().Format(time.RFC3339Nano)
	trackingPlantData := types.TrackingPlantData{
		SerNo: params.SerNo,
		IMEI: params.SerNo,
		ProdId: 1,
		FW: "1.1.1.1",
		Records: []types.TrackingPlantDataRecord{
			{
				SeqNo: params.SeqNo,
				Reason: 1,
				DateUTC: currentTime,
				Fields: []any{
					types.PlantGPSDataField{
						GpsUTC: currentTime,
						Lat: params.Latitude,
						Long: params.Longitude,
						Alt: 18,
						Spd: utils.ConvertKMHToMS(50),
						SpdAcc: 2,
						Head: 159,
						PDOP: 28,
						PosAcc: 6,
						GpsStat: 3,
						FType: 0,
					},
					types.PlantOdoField{
						Odo: 111,
						RH: 12345,
						FType: 1,
					},
					types.PlantDeviceField{
						DIn: 0,
						DOut: 0,
						DevStat: 2,
						FType: 2,
					},
					types.PlantAnalogueDataField{
						AnalogueData: map[int]int{
							4: utils.RandomNumberFromRange(15, 28),
							1: 4144,
							2: 30,
							5: 3,
							3: 1467,
						},
						FType: 1,
					},
				},
			},
		},
	}

	return trackingPlantData
}

func (client *Client) SendTrackingPlantData(params TrackingRequest) error {
	apiURL := client.config.GetAPIURL(params.State)
	apiKey := client.config.GetAPIKey(params.State)
	trackingPlantData := buildTrackingPlantData(BuildTrackingPlantDataBuilder{
		Latitude: params.Location.Latitude,
		Longitude: params.Location.Longitude,
		Speed: params.Location.Speed,
		SeqNo: params.SeqNo,
		SerNo: params.SerNo,
	})

	jsonData, err := json.Marshal(trackingPlantData)

	if err != nil {
		return fmt.Errorf("Error marshalling tracking plant data: %w", err)
	}

	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return fmt.Errorf("Error creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", apiKey)
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/114.0")
	response, err := client.httpClient.Do(request)

	if err != nil {
		return fmt.Errorf("Error sending tracking plant data: %w\n", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("Error reading response body: %w\n", err)
	}

	// consider another solution to convert this data
	fmt.Printf("Response-body: %s\n", string(body))
	fmt.Printf("Response-status: %d\n", response.StatusCode)

	return nil
}

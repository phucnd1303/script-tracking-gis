package api

import (
	"bytes"
	"context"
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
	config     *config.Config
	httpClient *http.Client
}

type TrackingRequest struct {
	State    string
	Location types.Location
	SeqNo    int
	SerNo    string
}

type APIResponse struct {
	Body       string
	StatusCode int
}

type BuildTrackingPlantDataBuilder struct {
	Latitude  string  `json:"lat"`
	Longitude string  `json:"log"`
	Speed     *string `json:"speed"`
	SeqNo     int
	SerNo     string
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}

func buildTrackingPlantData(params BuildTrackingPlantDataBuilder) (types.TrackingPlantData, error) {
	speed := "50"

	if params.Speed != nil {
		speed = *params.Speed
	}

	kmhSpeed, err := utils.ConvertKMHToCMS(speed)

	if err != nil {
		return types.TrackingPlantData{}, err
	}

	currentTime := time.Now().UTC().Format(time.RFC3339Nano)
	trackingPlantData := types.TrackingPlantData{
		SerNo:  params.SerNo,
		IMEI:   params.SerNo,
		ProdId: 1,
		FW:     "1.1.1.1",
		Records: []types.TrackingPlantDataRecord{
			{
				SeqNo:   params.SeqNo,
				Reason:  1,
				DateUTC: currentTime,
				Fields: []any{
					types.PlantGPSDataField{
						GpsUTC:  currentTime,
						Lat:     params.Latitude,
						Long:    params.Longitude,
						Alt:     18,
						Spd:     kmhSpeed,
						SpdAcc:  2,
						Head:    159,
						PDOP:    28,
						PosAcc:  6,
						GpsStat: 3,
						FType:   0,
					},
					types.PlantOdoField{
						Odo:   111,
						RH:    12345,
						FType: 1,
					},
					types.PlantDeviceField{
						DIn:     0,
						DOut:    0,
						DevStat: 2,
						FType:   2,
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

	return trackingPlantData, nil
}

func (client *Client) SendTrackingPlantData(ctx context.Context, params TrackingRequest) (*APIResponse, error) {
	apiURL := client.config.GetAPIURL(params.State)
	apiKey := client.config.GetAPIKey(params.State)
	trackingPlantData, err := buildTrackingPlantData(BuildTrackingPlantDataBuilder{
		Latitude:  params.Location.Latitude,
		Longitude: params.Location.Longitude,
		Speed:     params.Location.Speed,
		SeqNo:     params.SeqNo,
		SerNo:     params.SerNo,
	})

	if err != nil {
		return nil, fmt.Errorf("error building tracking plant data: %w", err)
	}

	jsonData, err := json.Marshal(trackingPlantData)

	if err != nil {
		return nil, fmt.Errorf("error marshalling tracking plant data: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", apiKey)
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/114.0")
	response, err := client.httpClient.Do(request)

	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("request cancelled: %w", ctx.Err())
		}

		return nil, fmt.Errorf("error sending tracking plant data: %w", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return &APIResponse{
		Body:       string(body),
		StatusCode: response.StatusCode,
	}, nil
}

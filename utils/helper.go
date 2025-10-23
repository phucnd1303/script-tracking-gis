package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func RandomNumberFromRange(min int, max int) int {
	result := rand.Intn(max - min + 1) + min

	return result
}

func ConvertKMHToMS(kmh int) float32 {
	result := float32(kmh * 1000) / 3600

	return result
}

func ConvertKMHToCMS(kmh int) float32 {
	result := float32(kmh) / 0.036

	return result
}

func GetAPIKey(state string) string {
	key := fmt.Sprintf("API_KEY_%s", strings.ToUpper(state))
	result := os.Getenv(key)

	return result
}

func GetPlantAPIURL(state string) string {
	url := fmt.Sprintf("API_URL_PLANT_%s", strings.ToUpper(state))
	result := os.Getenv(url)

	return result
}

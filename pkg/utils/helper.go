package utils

import "math/rand"

func RandomNumberFromRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func ConvertKMHToMS(kmh int) float32 {
	return float32(kmh*1000) / 3600
}

func ConvertKMHToCMS(kmh int) float32 {
	return float32(kmh) / 0.036
}

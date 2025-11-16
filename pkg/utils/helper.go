package utils

import (
	"math"
	"math/rand"
	"strconv"
)

func RandomNumberFromRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func ConvertKMHToMS(kmh string) (float64, error) {
	kmhValue, err := ParseFloat64(kmh)

	if err != nil {
		return 0, err
	}

	msValue := (kmhValue * 1000) / 3600
	roundedValue := RoundFloat(msValue, 2)

	return roundedValue, nil
}

func ConvertKMHToCMS(kmh string) (float64, error) {
	kmhValue, err := ParseFloat64(kmh)

	if err != nil {
		return 0, err
	}

	cmsValue := kmhValue / 0.036
	roundedValue := RoundFloat(cmsValue, 2)

	return roundedValue, nil
}

func ParseFloat64(s string) (float64, error) {
	value, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0, err
	}

	return float64(value), nil
}

func RoundFloat(value float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))

	return math.Round(value*ratio) / ratio
}

package models

type Location struct {
	Latitude  string  `json:"lat"`
	Longitude string  `json:"long"`
	Speed     *string `json:"speed"`
}

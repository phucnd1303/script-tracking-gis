package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKeys map[string]string
	APIURLs map[string]string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	return &Config{
		APIKeys: map[string]string {
			"local": os.Getenv("API_KEY_LOCAL"),
			"dev": os.Getenv("API_KEY_DEV"),
			"uat": os.Getenv("API_KEY_UAT"),
			"sit": os.Getenv("API_KEY_SIT"),
		},
		APIURLs: map[string]string {
			"local": os.Getenv("API_URL_PLANT_LOCAL"),
			"dev": os.Getenv("API_URL_PLANT_DEV"),
			"uat": os.Getenv("API_URL_PLANT_UAT"),
			"sit": os.Getenv("API_URL_PLANT_SIT"),
		},
	}, nil
}

func (config Config) GetAPIKey(env string) string {
	return config.APIKeys[env]
}

func (config Config) GetAPIURL(env string) string {
	return config.APIURLs[env]
}

package config

import (
	"errors"
	"os"
)

type App struct {
	RapidApiKey string
	S3Bucket    string
}

func Initialize() (App, error) {
	config := App{
		RapidApiKey: os.Getenv("RAPIDAPI_KEY"),
		S3Bucket:    os.Getenv("S3_BUCKET_NAME"),
	}

	if config.RapidApiKey == "" {
		return App{}, errors.New("RAPIDAPI_KEY is not set")
	}
	if config.S3Bucket == "" {
		return App{}, errors.New("S3_BUCKET_NAME is not set")
	}

	return config, nil
}

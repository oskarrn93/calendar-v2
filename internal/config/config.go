package config

import (
	"errors"
	"os"
)

type App struct {
	RapidApiKey          string
	RapidApiKeyParameter string
	S3Bucket             string
}

func Initialize() (App, error) {
	config := App{
		RapidApiKey:          os.Getenv("RAPIDAPI_KEY"),
		RapidApiKeyParameter: os.Getenv("RAPIDAPI_KEY_PARAMETER"),
		S3Bucket:             os.Getenv("S3_BUCKET_NAME"),
	}

	if config.S3Bucket == "" {
		return App{}, errors.New("S3_BUCKET_NAME is not set")
	}

	return config, nil
}

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type NBARapidApiConfig struct {
	baseUrl string
	season  int
}

type RapidApiConfig struct {
	nba    NBARapidApiConfig
	apiKey string
}

type AppConfig struct {
	rapidApi RapidApiConfig
	s3Bucket string
}

func InitializeConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file was provided")
	}

	rapidApiKey := os.Getenv("RAPIDAPI_KEY")
	s3Bucket := os.Getenv("S3_BUCKET_ARN")

	config := AppConfig{
		rapidApi: RapidApiConfig{
			nba: NBARapidApiConfig{
				baseUrl: "https://api-nba-v1.p.rapidapi.com",
				season:  2024,
			},
			apiKey: rapidApiKey,
		},
		s3Bucket: s3Bucket,
	}

	return config
}

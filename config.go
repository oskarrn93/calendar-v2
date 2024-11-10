package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	rapidApiKey string
	s3Bucket    string
}

func InitializeConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file was provided")
	}

	rapidApiKey := os.Getenv("RAPIDAPI_KEY")
	s3Bucket := os.Getenv("S3_BUCKET_ARN")

	config := AppConfig{rapidApiKey: rapidApiKey, s3Bucket: s3Bucket}
	return config
}

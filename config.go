package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	rapidApiKey string
}

func InitializeConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rapidApiKey := os.Getenv("RAPIDAPI_KEY")

	config := AppConfig{rapidApiKey: rapidApiKey}
	return config
}

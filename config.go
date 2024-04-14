package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	rapidApiKey string
}

func InitializeConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rapidApiKey := os.Getenv("RAPIDAPI_KEY")

	config := Config{rapidApiKey: rapidApiKey}
	return config
}

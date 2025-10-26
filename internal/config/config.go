package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type NBARapidApi struct {
	BaseUrl string
	Season  int
}

type RapidApi struct {
	NBA    NBARapidApi
	ApiKey string
}

type App struct {
	RapidApi RapidApi
	S3Bucket string
}

func Initialize(logger *slog.Logger) App {
	err := godotenv.Load()
	if err != nil {
		logger.Debug("No .env file was provided")
	}

	rapidApiKey := os.Getenv("RAPIDAPI_KEY")
	s3Bucket := os.Getenv("S3_BUCKET_ARN")

	config := App{
		RapidApi: RapidApi{
			NBA: NBARapidApi{
				BaseUrl: "https://api-nba-v1.p.rapidapi.com",
				Season:  2024,
			},
			ApiKey: rapidApiKey,
		},
		S3Bucket: s3Bucket,
	}

	return config
}

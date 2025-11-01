package config

import (
	"fmt"
	"log/slog"
	"os"

	validator "github.com/oskarrn93/calendar-v2/internal/validation"
)

type RapidApiResource struct {
	BaseUrl string `validate:"required"`
	Season  int    `validate:"required"`
}

type RapidApi struct {
	NBA      RapidApiResource `validate:"required"`
	Football RapidApiResource `validate:"required"`
	ApiKey   string           `validate:"required"`
}

type App struct {
	RapidApi RapidApi `validate:"required"`
	S3Bucket string   `validate:"required"`
}

func (a *App) Validate() error {
	return validator.ValidateStruct(a)
}

func Initialize(logger *slog.Logger) App {
	config := App{
		RapidApi: RapidApi{
			NBA: RapidApiResource{
				BaseUrl: "https://api-nba-v1.p.rapidapi.com",
				Season:  2025,
			},
			Football: RapidApiResource{
				BaseUrl: "https://api-football-v1.p.rapidapi.com",
				Season:  2025,
			},
			ApiKey: os.Getenv("RAPIDAPI_KEY"),
		},
		S3Bucket: os.Getenv("S3_BUCKET_NAME"),
	}

	if err := config.Validate(); err != nil {
		panic(fmt.Errorf("Config validation failed: %w", err))
	}

	return config
}

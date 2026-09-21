package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithPrivateFieldValidation())

type RapidApiResource struct {
	BaseUrl string `validate:"required"`
}

type RapidApi struct {
	NBA        RapidApiResource `validate:"required"`
	Football   RapidApiResource `validate:"required"`
	Basketball RapidApiResource `validate:"required"`
	Esport     RapidApiResource `validate:"required"`
	ApiKey     string           `validate:"required"`
}

type App struct {
	RapidApi RapidApi `validate:"required"`
	S3Bucket string   `validate:"required"`
}

func (a *App) Validate() error {
	return validate.Struct(a)
}

func Initialize() (App, error) {
	config := App{
		RapidApi: RapidApi{
			NBA: RapidApiResource{
				BaseUrl: "https://api-nba-v1.p.rapidapi.com",
			},
			Football: RapidApiResource{
				BaseUrl: "https://api-football-v1.p.rapidapi.com",
			},
			Basketball: RapidApiResource{
				BaseUrl: "https://sportapi7.p.rapidapi.com",
			},
			Esport: RapidApiResource{
				BaseUrl: "https://pinnacle-odds.p.rapidapi.com",
			},
			ApiKey: os.Getenv("RAPIDAPI_KEY"),
		},
		S3Bucket: os.Getenv("S3_BUCKET_NAME"),
	}

	if err := config.Validate(); err != nil {
		return App{}, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

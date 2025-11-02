package testutil

import "github.com/oskarrn93/calendar-v2/internal/config"

func GetMockAppConfig() config.App {
	config := config.App{
		RapidApi: config.RapidApi{
			NBA: config.RapidApiResource{
				BaseUrl: "https://example-nba.com",
			},
			Football: config.RapidApiResource{
				BaseUrl: "https://example-football.com",
			},
			Esport: config.RapidApiResource{
				BaseUrl: "https://example-esport.com",
			},
			ApiKey: "fake-api-key",
		},
		S3Bucket: "fake-s3-bucket",
	}

	if err := config.Validate(); err != nil {
		panic("Invalid mock app config")
	}

	return config
}

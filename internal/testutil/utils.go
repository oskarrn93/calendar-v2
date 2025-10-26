package testutil

import "oskarrn93/calendar-v2/internal/config"

func GetMockAppConfig() config.App {
	config := config.App{
		RapidApi: config.RapidApi{
			NBA: config.NBARapidApi{
				BaseUrl: "https://example.com",
				Season:  2024,
			},
			ApiKey: "fake-api-key",
		},
		S3Bucket: "fake-s3-bucket",
	}

	return config
}

package main

func getMockAppConfig() AppConfig {
	config := AppConfig{
		rapidApi: RapidApiConfig{
			nba: NBARapidApiConfig{
				baseUrl: "https://example.com",
				season:  2024,
			},
			apiKey: "fake-api-key",
		},
		s3Bucket: "fake-s3-bucket",
	}

	return config
}

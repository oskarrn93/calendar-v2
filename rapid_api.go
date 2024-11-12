package main

import (
	"log"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// Ap docs: https://rapidapi.com/api-sports/api/api-nba

type RapidApi struct {
	httpClient *resty.Client
	config     RapidApiConfig
}

func (ra RapidApi) getBaseRequest(baseUrl string) *resty.Request {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf("Faiiled to parse RapidApi base url: %s", baseUrl)
	}

	return ra.httpClient.R().EnableTrace().SetHeader("X-RapidAPI-Key", ra.config.apiKey).SetHeader("X-RapidAPI-Host", parsedUrl.Hostname())
}

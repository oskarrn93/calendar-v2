package main

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// Ap docs: https://rapidapi.com/api-sports/api/api-nba

type RapidApi struct {
	httpClient *resty.Client
	config     RapidApiConfig
}

func (ra RapidApi) getBaseRequest(baseUrl string) (*resty.Request, error) {
	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("Faiiled to parse RapidApi base url: %s: %w", baseUrl, err)
	}

	baseRequest := ra.httpClient.R().EnableTrace().SetHeader("X-RapidAPI-Key", ra.config.apiKey).SetHeader("X-RapidAPI-Host", parsedUrl.Hostname())

	return baseRequest, nil
}

package main

import (
	"github.com/go-resty/resty/v2"
)

// Ap docs: https://rapidapi.com/api-sports/api/api-nba

type RapidApi struct {
	httpClient *resty.Client
	config     RapidApiConfig
}

func (ra RapidApi) getBaseRequest() *resty.Request {
	return ra.httpClient.R().EnableTrace().SetHeader("X-RapidAPI-Key", ra.config.apiKey)
}

package main

import (
	"oskarrn93/calendar-v2/internal/config"

	"github.com/go-resty/resty/v2"
)

// Ap docs: https://rapidapi.com/api-sports/api/api-nba

type RapidApi struct {
	httpClient *resty.Client
	config     config.RapidApi
}

func (ra RapidApi) getBaseRequest() *resty.Request {
	return ra.httpClient.R().EnableTrace().SetHeader("X-RapidAPI-Key", ra.config.ApiKey)
}

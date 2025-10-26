package rapidapi

import (
	"github.com/oskarrn93/calendar-v2/internal/config"

	"github.com/go-resty/resty/v2"
)

// Ap docs: https://rapidapi.com/api-sports/api/api-nba

type RapidApi struct {
	HttpClient *resty.Client
	Config     config.RapidApi
}

func (ra RapidApi) BaseRequest() *resty.Request {
	return ra.HttpClient.R().EnableTrace().SetHeader("X-RapidAPI-Key", ra.Config.ApiKey)
}

func New(httpClient *resty.Client, config config.RapidApi) RapidApi {
	return RapidApi{
		HttpClient: httpClient,
		Config:     config,
	}
}

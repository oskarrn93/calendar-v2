package rapidapi

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/oskarrn93/calendar-v2/internal/config"
	"github.com/oskarrn93/calendar-v2/internal/httpclient"
)

type RapidApi struct {
	HttpClient *resty.Client
	Config     config.RapidApi
}

func (ra RapidApi) BaseRequest(ctx context.Context) *resty.Request {
	return ra.HttpClient.R().SetContext(ctx).SetHeader("X-RapidAPI-Key", ra.Config.ApiKey)
}

func New(httpClient *resty.Client, config config.RapidApi) RapidApi {
	return RapidApi{
		HttpClient: httpclient.Configure(httpClient),
		Config:     config,
	}
}

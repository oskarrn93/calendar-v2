package rapidapi

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type RapidApi struct {
	httpClient *resty.Client
	apiKey     string
}

func (ra RapidApi) BaseRequest(ctx context.Context) *resty.Request {
	return ra.httpClient.R().SetContext(ctx).SetHeader("X-RapidAPI-Key", ra.apiKey)
}

func New(httpClient *resty.Client, apiKey string) RapidApi {
	return RapidApi{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

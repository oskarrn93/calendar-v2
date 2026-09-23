package rapidapi

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/oskarrn93/calendar-v2/internal/config"
)

const requestTimeout = 10 * time.Second

type RapidApi struct {
	HttpClient *resty.Client
	Config     config.RapidApi
}

func (ra RapidApi) BaseRequest(ctx context.Context) *resty.Request {
	return ra.HttpClient.R().SetContext(ctx).SetHeader("X-RapidAPI-Key", ra.Config.ApiKey)
}

func New(httpClient *resty.Client, config config.RapidApi) RapidApi {
	httpClient.SetTimeout(requestTimeout)

	// resty doesn't treat non-2xx as an error. Without this, an error body (rate
	// limit, bad key) unmarshals into an empty response and the handler uploads
	// an empty calendar over the previous good one.
	httpClient.OnAfterResponse(func(_ *resty.Client, response *resty.Response) error {
		if response.IsError() {
			return fmt.Errorf("unexpected status %s from %s", response.Status(), response.Request.URL)
		}
		return nil
	})

	return RapidApi{
		HttpClient: httpClient,
		Config:     config,
	}
}

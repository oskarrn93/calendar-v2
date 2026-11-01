package httpclient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const requestTimeout = 10 * time.Second

func New() *resty.Client {
	httpClient := resty.New().SetTimeout(requestTimeout)

	// resty doesn't treat non-2xx as an error. Without this, an error body (rate
	// limit, bad key) unmarshals into an empty response and the handler uploads
	// an empty calendar over the previous good one.
	httpClient.OnAfterResponse(func(_ *resty.Client, response *resty.Response) error {
		if response.IsError() {
			return fmt.Errorf("unexpected status %s from %s", response.Status(), response.Request.URL)
		}
		return nil
	})

	return httpClient
}

// GetJSON unmarshals the body explicitly instead of using resty's SetResult,
// which only decodes when the response has a JSON content type.
func GetJSON[T any](request *resty.Request, url string) (T, error) {
	var data T

	response, err := request.Get(url)
	if err != nil {
		return data, fmt.Errorf("request failed: %w", err)
	}

	if err := json.Unmarshal(response.Body(), &data); err != nil {
		return data, fmt.Errorf("failed to unmarshal response from %s: %w", url, err)
	}

	return data, nil
}

package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	requestTimeout = 10 * time.Second
	retryCount     = 2
	retryWaitTime  = time.Second
	retryMaxWait   = 5 * time.Second
)

func New() *resty.Client {
	httpClient := resty.New().
		SetTimeout(requestTimeout).
		SetRetryCount(retryCount).
		SetRetryWaitTime(retryWaitTime).
		SetRetryMaxWaitTime(retryMaxWait).
		// Once a condition is set it alone decides whether to retry, which also
		// overrides resty's default of never retrying OnAfterResponse errors.
		AddRetryCondition(isTransient)

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

// isTransient retries network failures, rate limits and server errors, but not
// other 4xx such as a bad API key, which will fail the same way again.
func isTransient(response *resty.Response, err error) bool {
	if response == nil || response.RawResponse == nil {
		return err != nil
	}
	return response.StatusCode() == http.StatusTooManyRequests || response.StatusCode() >= http.StatusInternalServerError
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

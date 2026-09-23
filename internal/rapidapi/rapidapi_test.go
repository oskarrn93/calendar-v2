package rapidapi_test

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestBaseRequestFailsOnErrorStatus(t *testing.T) {
	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())
	defer httpmock.DeactivateAndReset()

	rapidApi := rapidapi.New(httpClient, testutil.GetMockAppConfig(t).RapidApi)

	httpmock.RegisterResponder("GET", "https://example.com/ok", httpmock.NewStringResponder(200, `{}`))
	httpmock.RegisterResponder("GET", "https://example.com/limited", httpmock.NewStringResponder(429, `{"message":"Too many requests"}`))

	_, err := rapidApi.BaseRequest(t.Context()).Get("https://example.com/ok")
	require.NoError(t, err)

	_, err = rapidApi.BaseRequest(t.Context()).Get("https://example.com/limited")
	require.ErrorContains(t, err, "429")
}

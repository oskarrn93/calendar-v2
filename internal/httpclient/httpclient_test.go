package httpclient_test

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewFailsOnErrorStatus(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	httpmock.RegisterResponder("GET", "https://example.com/ok", httpmock.NewStringResponder(200, `{}`))
	httpmock.RegisterResponder("GET", "https://example.com/limited", httpmock.NewStringResponder(429, `{"message":"Too many requests"}`))

	_, err := httpClient.R().Get("https://example.com/ok")
	require.NoError(t, err)

	_, err = httpClient.R().Get("https://example.com/limited")
	require.ErrorContains(t, err, "429")
}

func TestNewRetriesServerErrors(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	httpmock.RegisterResponder("GET", "https://example.com/flaky", httpmock.ResponderFromMultipleResponses(
		[]*http.Response{
			httpmock.NewStringResponse(503, `{}`),
			httpmock.NewStringResponse(200, `{}`),
		},
	))

	_, err := httpClient.R().Get("https://example.com/flaky")
	require.NoError(t, err)
	require.Equal(t, 2, httpmock.GetTotalCallCount())
}

func TestNewDoesNotRetryClientErrors(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	httpmock.RegisterResponder("GET", "https://example.com/unauthorized", httpmock.NewStringResponder(401, `{}`))

	_, err := httpClient.R().Get("https://example.com/unauthorized")
	require.ErrorContains(t, err, "401")
	require.Equal(t, 1, httpmock.GetTotalCallCount())
}

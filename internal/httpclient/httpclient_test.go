package httpclient_test

import (
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

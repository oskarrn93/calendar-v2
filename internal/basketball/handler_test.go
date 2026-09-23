package basketball_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/basketball"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetGames(t *testing.T) {
	// Arrange
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())

	mockConfig := testutil.GetMockAppConfig(t)

	rapidApi := rapidapi.New(httpClient, mockConfig.RapidApi)

	handler := basketball.NewHandler(rapidApi, testutil.NoopStorage{}, logging.New())

	gamesTestData := string(testutil.ReadTestData(t, "basketball/events.json"))

	// Mock http request

	pageUrl := func(page int) string {
		expectedUrl, err := url.Parse(mockConfig.RapidApi.Basketball.BaseUrl)
		require.NoError(t, err)
		expectedUrl.Path += fmt.Sprintf("/api/v1/team/%d/events/next/%d", basketball.RealMadrid, page)
		return expectedUrl.String()
	}

	firstPage := `{"events":[{"id":1,"tournament":{"name":"Liga ACB"},"homeTeam":{"name":"Real Madrid"},"awayTeam":{"name":"Barcelona"},"startTimestamp":1762200000}],"hasNextPage":true}`
	httpmock.RegisterResponder("GET", pageUrl(0), httpmock.NewStringResponder(200, firstPage))
	httpmock.RegisterResponder("GET", pageUrl(1), httpmock.NewStringResponder(200, gamesTestData))

	// Act
	result, err := handler.GetGames(t.Context(), []basketball.TeamID{basketball.RealMadrid})
	require.NoError(t, err)

	// Assert
	require.Equal(t, 2, httpmock.GetTotalCallCount())
	require.Equal(t, int64(1), result[0].ID)

	snaps.MatchSnapshot(t, result)
}

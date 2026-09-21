package football_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/football"
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

	handler := football.NewHandler(rapidApi, testutil.NoopStorage{}, logging.New())

	gamesTestData := string(testutil.ReadTestData(t, "football/fixtures/real_madrid.json"))

	// Mock http request

	expectedUrl, err := url.Parse(mockConfig.RapidApi.Football.BaseUrl)
	require.NoError(t, err)
	expectedUrl.Path += "/v3/fixtures"
	query := expectedUrl.Query()
	query.Set("season", fmt.Sprintf("%d", football.Season))
	query.Set("team", fmt.Sprintf("%d", football.RealMadrid))
	expectedUrl.RawQuery = query.Encode()

	httpmock.RegisterResponder("GET", expectedUrl.String(),
		httpmock.NewStringResponder(200, gamesTestData))

	// Act
	result, err := handler.GetGames([]football.TeamID{football.RealMadrid})
	require.NoError(t, err)

	// Assert
	for _, game := range result {
		if game.Team.Home.Id != int(football.RealMadrid) && game.Team.Away.Id != int(football.RealMadrid) {
			t.Errorf("Expected either home or away team to be REAL_MADRID_TEAM_ID, but got home: %d, away: %d", game.Team.Home.Id, game.Team.Away.Id)
		}
	}

	snaps.MatchSnapshot(t, result)
}

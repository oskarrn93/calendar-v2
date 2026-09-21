package nba_test

import (
	"fmt"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/nba"
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

	nbaHandler := nba.NewHandler(rapidApi, testutil.NoopStorage{}, logging.New())

	gamesTestData := string(testutil.ReadTestData(t, "nba/games/celtics.json"))

	// Mock http request
	expectedUrl := fmt.Sprintf("%s/games?season=%d&team=%d", mockConfig.RapidApi.NBA.BaseUrl, nba.Season, nba.BostonCeltics)
	httpmock.RegisterResponder("GET", expectedUrl,
		httpmock.NewStringResponder(200, gamesTestData))

	// Act
	result, err := nbaHandler.GetGames([]nba.TeamID{nba.BostonCeltics})
	require.NoError(t, err)

	// Assert
	for _, game := range result {
		if game.Teams.Home.Id != int(nba.BostonCeltics) && game.Teams.Visitors.Id != int(nba.BostonCeltics) {
			t.Errorf("Expected either home or visitor team to be CELTICS_TEAM_ID, but got home: %d, visitor: %d", game.Teams.Home.Id, game.Teams.Visitors.Id)
		}
	}

	snaps.MatchSnapshot(t, result)
}

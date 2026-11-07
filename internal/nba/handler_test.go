package nba_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/nba"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetGames(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	rapidApi := rapidapi.New(httpClient, testutil.FakeApiKey)

	nbaHandler := nba.NewHandler(rapidApi, &testutil.MemoryStorage{})

	expectedUrl := fmt.Sprintf("%s/games?season=%d&team=%d", nba.BaseUrl, nba.Season(time.Now()), nba.BostonCeltics)
	httpmock.RegisterResponder("GET", expectedUrl,
		httpmock.NewBytesResponder(200, testutil.ReadTestData(t, "nba/games/celtics.json")))

	result, err := nbaHandler.GetGames(t.Context(), []nba.TeamID{nba.BostonCeltics})
	require.NoError(t, err)

	for _, game := range result {
		if game.Teams.Home.Id != int(nba.BostonCeltics) && game.Teams.Visitors.Id != int(nba.BostonCeltics) {
			t.Errorf("Expected either home or visitor team to be CELTICS_TEAM_ID, but got home: %d, visitor: %d", game.Teams.Home.Id, game.Teams.Visitors.Id)
		}
	}

	snaps.MatchSnapshot(t, result)
}

func TestSeason(t *testing.T) {
	require.Equal(t, 2026, nba.Season(time.Date(2027, time.June, 30, 0, 0, 0, 0, time.UTC)))
	require.Equal(t, 2027, nba.Season(time.Date(2027, time.July, 1, 0, 0, 0, 0, time.UTC)))
}

package football_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/football"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetGames(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	rapidApi := rapidapi.New(httpClient, testutil.FakeApiKey)

	handler := football.NewHandler(rapidApi, &testutil.MemoryStorage{})

	expectedUrl := fmt.Sprintf("%s/v3/fixtures?season=%d&team=%d", football.BaseUrl, football.Season(football.RealMadrid, time.Now()), football.RealMadrid)
	httpmock.RegisterResponder("GET", expectedUrl,
		httpmock.NewBytesResponder(200, testutil.ReadTestData(t, "football/fixtures/real_madrid.json")))

	result, err := handler.GetGames(t.Context(), []football.TeamID{football.RealMadrid})
	require.NoError(t, err)

	for _, game := range result {
		if game.Team.Home.Id != int(football.RealMadrid) && game.Team.Away.Id != int(football.RealMadrid) {
			t.Errorf("Expected either home or away team to be REAL_MADRID_TEAM_ID, but got home: %d, away: %d", game.Team.Home.Id, game.Team.Away.Id)
		}
	}

	snaps.MatchSnapshot(t, result)
}

func TestSeason(t *testing.T) {
	tests := []struct {
		name string
		team football.TeamID
		now  time.Time
		want int
	}{
		{"club before summer is last year's season", football.RealMadrid, time.Date(2027, time.March, 1, 0, 0, 0, 0, time.UTC), 2026},
		{"club from July is this year's season", football.RealMadrid, time.Date(2027, time.July, 1, 0, 0, 0, 0, time.UTC), 2027},
		{"calendar year team", football.MalmoFF, time.Date(2027, time.March, 1, 0, 0, 0, 0, time.UTC), 2027},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, football.Season(tt.team, tt.now))
		})
	}
}

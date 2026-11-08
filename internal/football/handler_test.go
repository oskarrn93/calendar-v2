package football_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
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

	expectedUrl := fmt.Sprintf("%s/v3/fixtures?season=%d&team=%d", football.BaseUrl, football.Seasons(football.RealMadrid, time.Now())[0], football.RealMadrid)
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

func TestGetGamesFailsOnApiError(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	handler := football.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), &testutil.MemoryStorage{})

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200,
		`{"errors":{"requests":"You have reached the request limit for the day"},"response":[]}`))

	_, err := handler.GetGames(t.Context(), []football.TeamID{football.RealMadrid})
	require.ErrorContains(t, err, "request limit")
}

func TestHandlerUploadsCalendar(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	storage := &testutil.MemoryStorage{}
	handler := football.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), storage)

	body := testutil.ReadTestData(t, "football/fixtures/real_madrid.json")
	httpmock.RegisterNoResponder(httpmock.NewBytesResponder(200, body))

	var fixtures football.FixturesResponse
	require.NoError(t, json.Unmarshal(body, &fixtures))

	require.NoError(t, handler.Handler(t.Context()))

	cal, err := ics.ParseCalendar(bytes.NewReader(storage.Files["football.ics"]))
	require.NoError(t, err)

	// Every team and season is served the same fixtures, so this also covers
	// deduplication. The fixture file has one postponed game.
	events := cal.Events()
	require.Len(t, events, len(fixtures.Response)-1)

	want := fixtures.Response[0]
	require.Equal(t, fmt.Sprintf("football-%d", want.Fixture.Id), events[0].Id())
	require.Equal(t, want.Team.Home.Name+" - "+want.Team.Away.Name, events[0].GetProperty(ics.ComponentPropertySummary).Value)

	start, err := events[0].GetStartAt()
	require.NoError(t, err)
	require.True(t, start.Equal(time.Unix(want.Fixture.Timestamp, 0)))

	end, err := events[0].GetEndAt()
	require.NoError(t, err)
	require.Equal(t, 2*time.Hour, end.Sub(start))
}

func TestHandlerSkipsCancelledAndMarksTBD(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	storage := &testutil.MemoryStorage{}
	handler := football.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), storage)

	fixture := func(id int, status string) string {
		return fmt.Sprintf(`{"fixture":{"id":%d,"timestamp":1762200000,"status":{"short":"%s"}},"teams":{"home":{"name":"Home"},"away":{"name":"Away"}}}`, id, status)
	}
	body := fmt.Sprintf(`{"errors":[],"response":[%s,%s,%s,%s,%s]}`, fixture(1, "NS"), fixture(2, "TBD"), fixture(3, "PST"), fixture(4, "CANC"), fixture(5, "ABD"))
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, body))

	require.NoError(t, handler.Handler(t.Context()))

	cal, err := ics.ParseCalendar(bytes.NewReader(storage.Files["football.ics"]))
	require.NoError(t, err)

	events := cal.Events()
	require.Len(t, events, 2)
	require.Equal(t, "Home - Away", events[0].GetProperty(ics.ComponentPropertySummary).Value)
	require.Equal(t, "Home - Away (TBD)", events[1].GetProperty(ics.ComponentPropertySummary).Value)
}

func TestGetGamesFailsWhenOnlySomeTeamsAreEmpty(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	handler := football.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), &testutil.MemoryStorage{})

	url := fmt.Sprintf("%s/v3/fixtures?season=%d&team=%d", football.BaseUrl, football.Seasons(football.RealMadrid, time.Now())[0], football.RealMadrid)
	httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, testutil.ReadTestData(t, "football/fixtures/real_madrid.json")))
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"errors":[],"response":[]}`))

	_, err := handler.GetGames(t.Context(), []football.TeamID{football.RealMadrid, football.Arsenal})
	require.ErrorContains(t, err, fmt.Sprintf("team ids [%d]", football.Arsenal))
}

func TestGetGamesAllowsEveryTeamBeingEmpty(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	handler := football.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), &testutil.MemoryStorage{})

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"errors":[],"response":[]}`))

	games, err := handler.GetGames(t.Context(), []football.TeamID{football.RealMadrid, football.Arsenal})
	require.NoError(t, err)
	require.Empty(t, games)
}

func TestSeasons(t *testing.T) {
	tests := []struct {
		name string
		team football.TeamID
		now  time.Time
		want []int
	}{
		{"club before summer is last year's season", football.RealMadrid, time.Date(2027, time.March, 1, 0, 0, 0, 0, time.UTC), []int{2026}},
		{"club from July is this year's season", football.RealMadrid, time.Date(2027, time.July, 1, 0, 0, 0, 0, time.UTC), []int{2027}},
		{"calendar year team before summer spans both seasons", football.MalmoFF, time.Date(2027, time.January, 15, 0, 0, 0, 0, time.UTC), []int{2026, 2027}},
		{"calendar year team from July is this year's season", football.MalmoFF, time.Date(2027, time.July, 1, 0, 0, 0, 0, time.UTC), []int{2027}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, football.Seasons(tt.team, tt.now))
		})
	}
}

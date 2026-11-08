package nba_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
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

func TestHandlerUploadsCalendar(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	storage := &testutil.MemoryStorage{}
	handler := nba.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), storage)

	body := testutil.ReadTestData(t, "nba/games/celtics.json")
	httpmock.RegisterNoResponder(httpmock.NewBytesResponder(200, body))

	var games nba.GamesResponse
	require.NoError(t, json.Unmarshal(body, &games))

	require.NoError(t, handler.Handler(t.Context()))

	cal, err := ics.ParseCalendar(bytes.NewReader(storage.Files["nba.ics"]))
	require.NoError(t, err)

	// Every team is served the same games, so this also covers deduplication.
	events := cal.Events()
	require.Len(t, events, len(games.Response))

	want := games.Response[0]
	require.Equal(t, fmt.Sprintf("nba-%d", want.Id), events[0].Id())
	require.Equal(t, want.Teams.Home.Name+" - "+want.Teams.Visitors.Name, events[0].GetProperty(ics.ComponentPropertySummary).Value)

	start, err := events[0].GetStartAt()
	require.NoError(t, err)
	require.True(t, start.Equal(want.Date.Start))

	end, err := events[0].GetEndAt()
	require.NoError(t, err)
	require.Equal(t, 2*time.Hour, end.Sub(start))
}

func TestHandlerSkipsCancelledGames(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	storage := &testutil.MemoryStorage{}
	handler := nba.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), storage)

	game := func(id int, status string) string {
		return fmt.Sprintf(`{"id":%d,"date":{"start":"2026-11-03T00:00:00.000Z"},"status":{"long":"%s"},"teams":{"home":{"name":"Home"},"visitors":{"name":"Away"}}}`, id, status)
	}
	body := fmt.Sprintf(`{"errors":[],"response":[%s,%s,%s]}`, game(1, "Scheduled"), game(2, "Postponed"), game(3, "Canceled"))
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, body))

	require.NoError(t, handler.Handler(t.Context()))

	cal, err := ics.ParseCalendar(bytes.NewReader(storage.Files["nba.ics"]))
	require.NoError(t, err)

	events := cal.Events()
	require.Len(t, events, 1)
	require.Equal(t, "nba-1", events[0].Id())
}

func TestGetGamesFailsWhenOnlySomeTeamsAreEmpty(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	handler := nba.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), &testutil.MemoryStorage{})

	url := fmt.Sprintf("%s/games?season=%d&team=%d", nba.BaseUrl, nba.Season(time.Now()), nba.BostonCeltics)
	httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, testutil.ReadTestData(t, "nba/games/celtics.json")))
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"errors":[],"response":[]}`))

	_, err := handler.GetGames(t.Context(), []nba.TeamID{nba.BostonCeltics, nba.LosAngelesLakers})
	require.ErrorContains(t, err, fmt.Sprintf("team ids [%d]", nba.LosAngelesLakers))
}

func TestGetGamesAllowsEveryTeamBeingEmpty(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	handler := nba.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), &testutil.MemoryStorage{})

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, `{"errors":[],"response":[]}`))

	games, err := handler.GetGames(t.Context(), []nba.TeamID{nba.BostonCeltics, nba.LosAngelesLakers})
	require.NoError(t, err)
	require.Empty(t, games)
}

func TestGetGamesFailsOnApiError(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	nbaHandler := nba.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), &testutil.MemoryStorage{})

	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200,
		`{"errors":{"token":"Error/Missing application key"},"response":[]}`))

	_, err := nbaHandler.GetGames(t.Context(), []nba.TeamID{nba.BostonCeltics})
	require.ErrorContains(t, err, "application key")
}

func TestSeason(t *testing.T) {
	require.Equal(t, 2026, nba.Season(time.Date(2027, time.June, 30, 0, 0, 0, 0, time.UTC)))
	require.Equal(t, 2027, nba.Season(time.Date(2027, time.July, 1, 0, 0, 0, 0, time.UTC)))
}

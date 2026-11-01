package esport_test

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"testing"

	ics "github.com/arran4/golang-ical"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/esport"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestGetGames(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	rapidApi := rapidapi.New(httpClient, testutil.FakeApiKey)

	handler := esport.NewHandler(rapidApi, &testutil.MemoryStorage{}, logging.New())

	expectedUrl := fmt.Sprintf("%s/kit/v1/markets?sport_id=%d", esport.BaseUrl, esport.EsportSportID)
	httpmock.RegisterResponder("GET", expectedUrl,
		httpmock.NewBytesResponder(200, testutil.ReadTestData(t, "esport/markets.json")))

	result, err := handler.GetEvents(t.Context())
	require.NoError(t, err)

	snaps.MatchSnapshot(t, result)
}

func TestHandlerUploadsFilteredCalendar(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	storage := &testutil.MemoryStorage{}
	handler := esport.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), storage, logging.New())

	httpmock.RegisterResponder("GET", esport.BaseUrl+"/kit/v1/markets",
		httpmock.NewBytesResponder(200, testutil.ReadTestData(t, "esport/markets.json")))

	require.NoError(t, handler.Handler(t.Context()))

	cal, err := ics.ParseCalendar(bytes.NewReader(storage.Files["esport.ics"]))
	require.NoError(t, err)

	events := cal.Events()
	require.NotEmpty(t, events)

	for _, event := range events {
		start, err := event.GetStartAt()
		require.NoError(t, err)
		end, err := event.GetEndAt()
		require.NoError(t, err)
		require.True(t, end.After(start), "event %s ends before it starts", event.Id())

		summary := event.GetProperty(ics.ComponentPropertySummary).Value
		require.True(t, slices.ContainsFunc(esport.TeamsOfInterest, func(team string) bool {
			return strings.Contains(strings.ToLower(summary), strings.ToLower(team))
		}), "unexpected event %q", summary)
	}
}

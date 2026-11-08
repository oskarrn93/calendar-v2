package basketball_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/basketball"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestHandlerUploadsCalendar(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	storage := &testutil.MemoryStorage{}
	handler := basketball.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), storage)

	body := `{"events":[{"id":1,"tournament":{"name":"Liga ACB"},"homeTeam":{"name":"Real Madrid","shortName":"R. Madrid"},"awayTeam":{"name":"Barcelona"},"startTimestamp":1762200000}],"hasNextPage":false}`
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, body))

	require.NoError(t, handler.Handler(t.Context()))

	cal, err := ics.ParseCalendar(bytes.NewReader(storage.Files["basketball.ics"]))
	require.NoError(t, err)

	events := cal.Events()
	require.Len(t, events, 1)
	require.Equal(t, "basketball-1", events[0].Id())
	require.Equal(t, "Liga ACB: R. Madrid - Barcelona", events[0].GetProperty(ics.ComponentPropertySummary).Value)

	start, err := events[0].GetStartAt()
	require.NoError(t, err)
	require.True(t, start.Equal(time.Unix(1762200000, 0)))

	end, err := events[0].GetEndAt()
	require.NoError(t, err)
	require.Equal(t, 2*time.Hour, end.Sub(start))
}

func TestHandlerSkipsCancelledGames(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	storage := &testutil.MemoryStorage{}
	handler := basketball.NewHandler(rapidapi.New(httpClient, testutil.FakeApiKey), storage)

	body := `{"events":[{"id":1,"status":{"type":"notstarted"},"startTimestamp":1762200000},{"id":2,"status":{"type":"postponed"},"startTimestamp":1762200000},{"id":3,"status":{"type":"canceled"},"startTimestamp":1762200000}],"hasNextPage":false}`
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, body))

	require.NoError(t, handler.Handler(t.Context()))

	cal, err := ics.ParseCalendar(bytes.NewReader(storage.Files["basketball.ics"]))
	require.NoError(t, err)

	events := cal.Events()
	require.Len(t, events, 1)
	require.Equal(t, "basketball-1", events[0].Id())
}

func TestGetGames(t *testing.T) {
	httpClient := testutil.NewMockHttpClient(t)

	rapidApi := rapidapi.New(httpClient, testutil.FakeApiKey)

	handler := basketball.NewHandler(rapidApi, &testutil.MemoryStorage{})

	pageUrl := func(page int) string {
		return fmt.Sprintf("%s/api/v1/team/%d/events/next/%d", basketball.BaseUrl, basketball.RealMadrid, page)
	}

	firstPage := `{"events":[{"id":1,"tournament":{"name":"Liga ACB"},"homeTeam":{"name":"Real Madrid"},"awayTeam":{"name":"Barcelona"},"startTimestamp":1762200000}],"hasNextPage":true}`
	httpmock.RegisterResponder("GET", pageUrl(0), httpmock.NewStringResponder(200, firstPage))
	httpmock.RegisterResponder("GET", pageUrl(1), httpmock.NewBytesResponder(200, testutil.ReadTestData(t, "basketball/events.json")))

	result, err := handler.GetGames(t.Context(), []basketball.TeamID{basketball.RealMadrid})
	require.NoError(t, err)

	require.Equal(t, 2, httpmock.GetTotalCallCount())
	require.Equal(t, int64(1), result[0].ID)

	snaps.MatchSnapshot(t, result)
}

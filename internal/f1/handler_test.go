package f1_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/f1"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

const unpublishedSeason = `{"MRData":{"RaceTable":{"Races":[]}}}`

// setup serves the given body for each season; seasons without a body have no
// responder, so httpmock returns an error and an unexpected request fails the test.
func setup(t *testing.T, bodies map[int][]byte) (*f1.Handler, *testutil.MemoryStorage) {
	t.Helper()

	httpClient := testutil.NewMockHttpClient(t)

	for season, body := range bodies {
		url := fmt.Sprintf("%s/ergast/f1/%d/races/?limit=100", f1.BaseUrl, season)
		httpmock.RegisterResponder("GET", url, httpmock.NewBytesResponder(200, body))
	}

	storage := &testutil.MemoryStorage{}
	return f1.NewHandler(httpClient, storage, logging.New()), storage
}

func TestGetRaces(t *testing.T) {
	handler, _ := setup(t, map[int][]byte{2026: testutil.ReadTestData(t, "f1/races.json")})

	result, err := handler.GetRaces(t.Context(), []int{2026})
	require.NoError(t, err)

	snaps.MatchSnapshot(t, result)
}

func TestHandlerFetchesCurrentAndNextSeason(t *testing.T) {
	year := time.Now().UTC().Year()
	handler, storage := setup(t, map[int][]byte{
		year:     testutil.ReadTestData(t, "f1/races.json"),
		year + 1: []byte(unpublishedSeason),
	})

	require.NoError(t, handler.Handler(t.Context()))
	require.Equal(t, 2, httpmock.GetTotalCallCount())

	cal, err := ics.ParseCalendar(bytes.NewReader(storage.Files["f1.ics"]))
	require.NoError(t, err)

	// 23 races and qualifyings plus 6 sprint weekends with sprint and sprint qualifying.
	require.Len(t, cal.Events(), 23*2+6*2)
}

func TestHandlerDoesNotUploadEmptySchedule(t *testing.T) {
	year := time.Now().UTC().Year()
	handler, storage := setup(t, map[int][]byte{
		year:     []byte(unpublishedSeason),
		year + 1: []byte(unpublishedSeason),
	})

	require.ErrorIs(t, handler.Handler(t.Context()), calendar.ErrNoEvents)
	require.Empty(t, storage.Files)
}

func TestHandlerDoesNotUploadWhenNextSeasonFails(t *testing.T) {
	year := time.Now().UTC().Year()
	handler, storage := setup(t, map[int][]byte{year: testutil.ReadTestData(t, "f1/races.json")})

	require.Error(t, handler.Handler(t.Context()))
	require.Empty(t, storage.Files)
}

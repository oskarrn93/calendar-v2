package esport_test

import (
	"bytes"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"testing"

	ics "github.com/arran4/golang-ical"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/esport"
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

	handler := esport.NewHandler(rapidApi, testutil.NoopStorage{}, logging.New())

	gamesTestData := string(testutil.ReadTestData(t, "esport/markets.json"))

	// Mock http request

	expectedUrl, err := url.Parse(mockConfig.RapidApi.Esport.BaseUrl)
	require.NoError(t, err)
	expectedUrl.Path += "/kit/v1/markets"
	query := expectedUrl.Query()
	query.Set("sport_id", fmt.Sprintf("%d", esport.EsportSportID))
	expectedUrl.RawQuery = query.Encode()

	httpmock.RegisterResponder("GET", expectedUrl.String(),
		httpmock.NewStringResponder(200, gamesTestData))

	// Act
	result, err := handler.GetEvents(t.Context(), []esport.SportID{esport.EsportSportID})
	require.NoError(t, err)

	// Assert

	snaps.MatchSnapshot(t, result)
}

func TestHandlerUploadsFilteredCalendar(t *testing.T) {
	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())
	defer httpmock.DeactivateAndReset()

	mockConfig := testutil.GetMockAppConfig(t)
	storage := &testutil.MemoryStorage{}
	handler := esport.NewHandler(rapidapi.New(httpClient, mockConfig.RapidApi), storage, logging.New())

	httpmock.RegisterResponder("GET", mockConfig.RapidApi.Esport.BaseUrl+"/kit/v1/markets",
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

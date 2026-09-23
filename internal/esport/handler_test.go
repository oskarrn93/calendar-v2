package esport_test

import (
	"fmt"
	"net/url"
	"testing"

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

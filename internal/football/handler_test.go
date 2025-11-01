package football_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/url"
	"testing"
	"time"

	"github.com/oskarrn93/calendar-v2/internal/football"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testdata"
	"github.com/oskarrn93/calendar-v2/internal/testutil"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func readGamesTestData(t *testing.T, teamID football.TeamID) []byte {
	// Use saved api response so we don't need to make an external request

	var filePath string
	switch teamID {
	case football.REAL_MADRID_TEAM_ID:
		filePath = "football/fixtures/real_madrid.json"
	case football.MALMO_FF_TEAM_ID:
		filePath = "football/fixtures/malmo_ff.json"
	default:
		t.Fatalf("No test data for team ID: %d", teamID)
	}

	jsonFile, err := testdata.Content.Open(filePath)
	if err != nil {
		log.Fatal("Failed to open games test data", err)
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Failed to read games test data", err)
	}

	return data
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Upload(ctx context.Context, filename string, data []byte, logger *slog.Logger) error {
	args := m.Called(ctx, filename, data, logger)
	return args.Error(0)
}

func TestGetGames(t *testing.T) {
	//Arrange
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())

	mockConfig := testutil.GetMockAppConfig()

	rapidApi := rapidapi.New(httpClient, mockConfig.RapidApi)
	mockStorage := MockStorage{}

	handler := football.NewHandler(rapidApi, &mockStorage, logging.New())

	gamesTestData := string(readGamesTestData(t, football.REAL_MADRID_TEAM_ID))

	expectedFromDate := time.Now().UTC().Add(-1 * time.Hour * 24 * 2).Format("2006-01-02")
	// Mock http request

	expectedUrl, err := url.Parse(mockConfig.RapidApi.Football.BaseUrl)
	require.NoError(t, err)
	expectedUrl.Path += "/fixtures"
	query := expectedUrl.Query()
	query.Set("from", expectedFromDate)
	query.Set("season", fmt.Sprintf("%d", mockConfig.RapidApi.Football.Season))
	query.Set("team", fmt.Sprintf("%d", football.REAL_MADRID_TEAM_ID))
	expectedUrl.RawQuery = query.Encode()

	httpmock.RegisterResponder("GET", expectedUrl.String(),
		httpmock.NewStringResponder(200, gamesTestData))

	// Act
	result, err := handler.GetGames([]football.TeamID{football.REAL_MADRID_TEAM_ID})
	require.NoError(t, err)

	// Assert
	for _, game := range result {
		if game.Team.Home.Id != int(football.REAL_MADRID_TEAM_ID) && game.Team.Away.Id != int(football.REAL_MADRID_TEAM_ID) {
			t.Errorf("Expected either home or away team to be REAL_MADRID_TEAM_ID, but got home: %d, away: %d", game.Team.Home.Id, game.Team.Away.Id)
		}
	}

	snaps.MatchSnapshot(t, result)
}

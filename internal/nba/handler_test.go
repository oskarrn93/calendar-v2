package nba_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/oskarrn93/calendar-v2/internal/logging"
	"github.com/oskarrn93/calendar-v2/internal/nba"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"github.com/oskarrn93/calendar-v2/internal/testdata"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func readGamesTestData(t *testing.T) []byte {
	// Use saved api response so we don't need to make an external request

	jsonFile, err := testdata.Content.Open("nba/games/celtics.json")
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to open games test data: %w", err))
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed to read games test data: %w", err))
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
	// Arrange
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())

	mockConfig := testutil.GetMockAppConfig()

	rapidApi := rapidapi.New(httpClient, mockConfig.RapidApi)
	mockStorage := MockStorage{}

	nbaHandler := nba.NewHandler(rapidApi, &mockStorage, logging.New())

	gamesTestData := string(readGamesTestData(t))

	// Mock http request
	expectedUrl := fmt.Sprintf("%s/games?season=%d&team=%d", mockConfig.RapidApi.NBA.BaseUrl, mockConfig.RapidApi.NBA.Season, nba.CELTICS_TEAM_ID)
	httpmock.RegisterResponder("GET", expectedUrl,
		httpmock.NewStringResponder(200, gamesTestData))

	// Act
	result, err := nbaHandler.GetGames([]nba.TeamID{nba.CELTICS_TEAM_ID})
	require.NoError(t, err)

	// Assert
	for _, game := range result {
		if game.Teams.Home.Id != int(nba.CELTICS_TEAM_ID) && game.Teams.Visitors.Id != int(nba.CELTICS_TEAM_ID) {
			t.Errorf("Expected either home or visitor team to be CELTICS_TEAM_ID, but got home: %d, visitor: %d", game.Teams.Home.Id, game.Teams.Visitors.Id)
		}
	}

	snaps.MatchSnapshot(t, result)
}

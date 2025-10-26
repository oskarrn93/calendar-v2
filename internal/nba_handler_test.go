package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"oskarrn93/calendar-v2/internal/testdata"
	"oskarrn93/calendar-v2/internal/testutil"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readGamesTestData() []byte {
	// Use saved api response so we don't need to make an external request

	jsonFile, err := testdata.Content.Open("nba/games/celtics.json")
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

type MockStorage struct{}

func (m *MockStorage) upload(ctx context.Context, filename string, data []byte) error {
	return nil
}

func TestGetGames(t *testing.T) {
	//Arrange
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())

	mockConfig := testutil.GetMockAppConfig()

	rapidApi := RapidApi{httpClient: httpClient, config: mockConfig.RapidApi}
	mockStorage := MockStorage{}

	nbaHandler := NbaHandler{
		rapidApi: rapidApi,
		storage:  &mockStorage,
		logger:   NewLogger(),
	}

	gamesTestData := string(readGamesTestData())

	// Mock http request
	expectedUrl := fmt.Sprintf("%s/games?season=%d&team=%d", mockConfig.RapidApi.NBA.BaseUrl, mockConfig.RapidApi.NBA.Season, CELTICS_TEAM_ID)
	httpmock.RegisterResponder("GET", expectedUrl,
		httpmock.NewStringResponder(200, gamesTestData))

	// Act
	result, err := nbaHandler.getGames([]NBATeamID{CELTICS_TEAM_ID})
	require.NoError(t, err)

	// Assert
	assert.Len(t, result, 87)

	for _, game := range result {
		if game.Teams.Home.Id != int(CELTICS_TEAM_ID) && game.Teams.Visitors.Id != int(CELTICS_TEAM_ID) {
			t.Errorf("Expected either home or visitor team to be CELTICS_TEAM_ID, but got home: %d, visitor: %d", game.Teams.Home.Id, game.Teams.Visitors.Id)
		}
	}
}

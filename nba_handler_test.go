package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

func readGamesTestData() []byte {
	// Use saved api response so we don't need to make an external request

	jsonFile, err := os.Open("tests/data/nba/games/celtics.json")
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

	mockConfig := getMockAppConfig()

	rapidApi := RapidApi{httpClient: httpClient, config: mockConfig.rapidApi}
	mockStorage := MockStorage{}

	nbaHandler := NbaHandler{
		rapidApi: rapidApi,
		storage:  &mockStorage,
		logger:   NewLogger(),
	}

	// Mock http request
	expectedUrl := fmt.Sprintf("%s/games?season=%d&team=%d", mockConfig.rapidApi.nba.baseUrl, mockConfig.rapidApi.nba.season, nbaTeams.Celtics)
	httpmock.RegisterResponder("GET", expectedUrl,
		httpmock.NewStringResponder(200, string(readGamesTestData())))

	// Act
	result, err := nbaHandler.getGames()
	if err != nil {
		t.Fatalf("Did not expect error when calling nbaHandler.getGames: %v", err)
	}

	// Assert
	expectedTestData, err := nbaHandler.parseGamesResponse(readGamesTestData())
	if err != nil {
		t.Fatalf("Did not expect error when calling nbaHandler.parseGamesResponse: %v", err)
	}

	if reflect.DeepEqual(result, expectedTestData.Response) == false {
		t.Fatalf("Does not match expected data")
	}

}

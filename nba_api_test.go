package main

import (
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

func TestParseGames(t *testing.T) {
	//Arrange
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())

	fakeApikey := "fake-api-key"

	nbaApi := NBAApi{httpClient: httpClient, apiKey: fakeApikey}

	// Act
	nbaGames, err := nbaApi.parseGames(readGamesTestData())
	if err != nil {
		log.Fatal("Failed to parse games", err)
	}

	// Assert
	expectedNumberOfGames := 87 //from hard coded test data

	if nbaGames.Results != expectedNumberOfGames {
		log.Fatalf("Not expected number of games")
	}
	if len(nbaGames.Response) != expectedNumberOfGames {
		log.Fatalf("Response array is not matching expected number of games")
	}

	expectedTeamName := "Boston Celtics"
	for i := range nbaGames.Response {
		game := nbaGames.Response[i]

		if expectedTeamName != game.Teams.Home.Name && expectedTeamName != game.Teams.Visitors.Name {
			log.Fatalf("None of the teams is the expected team")
		}
	}
}

func TestGetGames(t *testing.T) {
	//Arrange
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpClient := resty.New()
	httpmock.ActivateNonDefault(httpClient.GetClient())

	// Mock http request
	expectedUrl := fmt.Sprintf("%s/games?season=%d&team=%d", NBABaseUrl, NBASeason, NBATeamIds["celtics"])
	httpmock.RegisterResponder("GET", expectedUrl,
		httpmock.NewStringResponder(200, string(readGamesTestData())))

	fakeApikey := "fake-api-key"
	nbaApi := NBAApi{httpClient: httpClient, apiKey: fakeApikey}

	// Act
	result := nbaApi.getGames()

	// Assert
	expectedResult, _ := nbaApi.parseGames(readGamesTestData())

	if reflect.DeepEqual(result, expectedResult) == false {
		log.Fatalf("Does not match expected data")
	}

}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// Ap docs: https://rapidapi.com/api-sports/api/api-nba

var NBATeamIds = map[string]int{
	"celtics": 2,
}
var NBASeason = 2023

var NBABaseUrl = "https://api-nba-v1.p.rapidapi.com"

type NBAApi struct {
	httpClient resty.Client
	appConfig  AppConfig
}

func (api NBAApi) getBaseRequest() *resty.Request {
	baseUrl, err := url.Parse(NBABaseUrl)
	if err != nil {
		log.Fatal("Faiiled to parse NBA Api base url")
	}

	return api.httpClient.R().EnableTrace().SetHeader("X-RapidAPI-Key", api.appConfig.rapidApiKey).SetHeader("X-RapidAPI-Host", baseUrl.Hostname())
}

func (api NBAApi) getGames() NBAGamesResponse {

	queryParams := map[string]string{
		"team":   strconv.Itoa(NBATeamIds["celtics"]),
		"season": strconv.Itoa(NBASeason),
	}

	apiUrl, err := url.Parse(fmt.Sprintf("%s/games", NBABaseUrl))
	if err != nil {
		log.Fatal("Faiiled to parse NBA Api games url")
	}

	request := api.getBaseRequest()
	response, err := request.SetQueryParams(queryParams).Get(apiUrl.String())
	if err != nil {
		log.Fatal("Request failed to retrieve NBA games", err)
	}

	log.Printf("response: %v", response)

	result, err := api.parseGames(response.Body())
	if err != nil {
		log.Fatal("Failed to parse games", err)
	}

	return result
}

type NBAGameDate struct {
	Start time.Time `json:time`
}

type NBATeam struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type NBAGameTeams struct {
	Visitors NBATeam `json:"visitors"`
	Home     NBATeam `json:"home"`
}

type NBAGame struct {
	Id    int          `json:"id"`
	Date  NBAGameDate  `json:"date"`
	Teams NBAGameTeams `json:"teams"`
}

type NBAGamesResponse struct {
	Get      string    `json:"get"`
	Results  int       `json:"results"`
	Response []NBAGame `json:"response"`
}

func (api NBAApi) parseGames(input []byte) (NBAGamesResponse, error) {
	var data NBAGamesResponse
	err := json.Unmarshal(input, &data)
	if err != nil {
		return data, errors.New("failed to unmarshall NBA games")
	}

	return data, nil

}

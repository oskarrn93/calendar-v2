package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"
)

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

type NBATeams struct {
	Celtics int
}

var nbaTeams = NBATeams{
	Celtics: 2,
}

type NbaHandler struct {
	rapidApi RapidApi
	storage  Storage
	logger   *slog.Logger
}

func (n *NbaHandler) handler(ctx context.Context) error {
	games, err := n.getGames()
	if err != nil {
		return err
	}

	calendar := n.createCalendar(games)
	calendarData := calendar.Export()

	if err := n.storage.upload(ctx, "nba.ics", calendarData); err != nil {
		return fmt.Errorf("failed to upload NBA file: %w", err)
	}

	return nil
}

func (n *NbaHandler) getGamesByTeam(teamId int) (NBAGamesResponse, error) {

	// TODO: Add support for multiple teams
	queryParams := map[string]string{
		"team":   strconv.Itoa(teamId),
		"season": strconv.Itoa(n.rapidApi.config.nba.season),
	}

	apiUrl, err := url.Parse(fmt.Sprintf("%s/games", n.rapidApi.config.nba.baseUrl))
	if err != nil {
		return NBAGamesResponse{}, fmt.Errorf("Faiiled to parse NBA Api games url: %w", err)
	}

	response, err := n.rapidApi.getBaseRequest().SetQueryParams(queryParams).Get(apiUrl.String())
	if err != nil {
		return NBAGamesResponse{}, fmt.Errorf("Request failed to retrieve NBA games: %w", err)
	}

	n.logger.Debug("NBA Api response", "response", response)

	return n.parseGamesResponse(response.Body())

}

func (n *NbaHandler) parseGamesResponse(input []byte) (NBAGamesResponse, error) {
	var data NBAGamesResponse
	if err := json.Unmarshal(input, &data); err != nil {
		return data, fmt.Errorf("Failed to unmarshall nba games: %w", err)
	}

	return data, nil
}

func (n *NbaHandler) getGames() ([]NBAGame, error) {
	var data = make([]NBAGame, 0)

	celticsData, err := n.getGamesByTeam(nbaTeams.Celtics)
	if err != nil {
		return data, err
	}

	data = append(data, celticsData.Response...)

	return data, nil

}

func (n *NbaHandler) createCalendar(games []NBAGame) Calendar {
	calendar := NewCalendar("NBA")
	for _, game := range games {
		newEvent := CalendarEvent{
			Id:        fmt.Sprintf("nba-%d", game.Id),
			Title:     fmt.Sprintf("%s - %s", game.Teams.Home.Name, game.Teams.Visitors.Name),
			StartDate: game.Date.Start,
			EndDate:   game.Date.Start.Add(2 * time.Hour),
		}
		calendar.AddEvent(newEvent)
	}

	return calendar
}

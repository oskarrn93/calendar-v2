package football

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
)

type Handler struct {
	rapidApi rapidapi.RapidApi
	storage  awsutil.Storage
	logger   *slog.Logger
}

func (h *Handler) Handler(ctx context.Context) error {
	games, err := h.GetGames(TeamIDs)
	if err != nil {
		return err
	}

	h.logger.Debug("Fotball games", "games", games)

	calendar := h.createCalendar(games)
	calendarData := calendar.Export()

	if err := h.storage.Upload(ctx, "football.ics", calendarData, h.logger); err != nil {
		return fmt.Errorf("failed to upload Football file: %w", err)
	}

	return nil
}

type Team struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type FixtureTeam struct {
	Home Team `json:"home"`
	Away Team `json:"away"`
}

type Fixture struct {
	Fixture struct {
		Id        int64 `json:"id"`
		Timestamp int64 `json:"timestamp"`
	} `json:"fixture"`
	Team FixtureTeam `json:"teams"`
}

type FixturesResponse struct {
	Get      string    `json:"get"`
	Results  int       `json:"results"`
	Response []Fixture `json:"response"`
}

func (f *Handler) GetGames(teamIDs []TeamID) ([]Fixture, error) {
	games := []Fixture{}

	for _, teamID := range teamIDs {
		response, err := f.getGamesByTeam(int(teamID))
		if err != nil {
			return nil, err
		}

		games = append(games, response.Response...)
	}

	return games, nil
}

func (f *Handler) getGamesByTeam(teamId int) (FixturesResponse, error) {
	/*
		curl -X GET https://api-football-v1.p.rapidapi.com/v3/fixtures?team=541&season=2024&from=2024-11-12 \
			--header 'x-rapidapi-key: REPLACE_ME' | jq .
	*/

	// TODO: Add support for multiple teams
	queryParams := map[string]string{
		"team":   strconv.Itoa(teamId),
		"season": strconv.Itoa(f.rapidApi.Config.Football.Season),
		"from":   time.Now().UTC().Add(-1 * time.Hour * 24 * 2).Format("2006-01-02"),
	}

	apiUrl, err := url.Parse(fmt.Sprintf("%s/fixtures", f.rapidApi.Config.Football.BaseUrl))
	if err != nil {
		return FixturesResponse{}, fmt.Errorf("Faiiled to parse Football Api games url: %w", err)
	}

	response, err := f.rapidApi.BaseRequest().SetQueryParams(queryParams).Get(apiUrl.String())
	if err != nil {
		return FixturesResponse{}, fmt.Errorf("Request failed to retrieve Football games: %w", err)
	}

	f.logger.Debug("Fotball Api response", "response", response)

	return f.parseGamesResponse(response.Body())
}

func (f *Handler) parseGamesResponse(input []byte) (FixturesResponse, error) {
	var data FixturesResponse
	if err := json.Unmarshal(input, &data); err != nil {
		return data, fmt.Errorf("Failed to unmarshall Football games: %w", err)
	}

	return data, nil
}

func (f *Handler) createCalendar(games []Fixture) calendar.Calendar {
	cal := calendar.New("Football")
	for _, game := range games {

		startTime := time.Unix(game.Fixture.Timestamp, 0)

		newEvent := calendar.Event{
			Id:        fmt.Sprintf("football-%d", game.Fixture.Id),
			Title:     fmt.Sprintf("%s - %s", game.Team.Home.Name, game.Team.Away.Name),
			StartDate: startTime,
			EndDate:   startTime.Add(2 * time.Hour),
		}
		cal.AddEvent(newEvent)
	}

	return cal
}

func NewHandler(rapidApi rapidapi.RapidApi, storage awsutil.Storage, logger *slog.Logger) *Handler {
	return &Handler{
		rapidApi: rapidApi,
		storage:  storage,
		logger:   logger,
	}
}

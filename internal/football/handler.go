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
	calendarData, err := calendar.Export()
	if err != nil {
		return fmt.Errorf("failed to export Football calendar: %w", err)
	}

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

func (h *Handler) GetGames(teamIDs []TeamID) ([]Fixture, error) {
	games := []Fixture{}

	for _, teamID := range teamIDs {
		response, err := h.getGamesByTeam(int(teamID))
		if err != nil {
			return nil, err
		}

		h.logger.Debug("Retrieved football games", "teamId", teamID, "results", response.Results)
		games = append(games, response.Response...)
	}

	return games, nil
}

func (h *Handler) getGamesByTeam(teamId int) (FixturesResponse, error) {
	/*
		curl -X GET https://api-football-v1.p.rapidapi.com/v3/fixtures?team=541&season=2025 \
			--header 'x-rapidapi-key: REPLACE_ME' | jq .
	*/

	// TODO: Add support for multiple teams
	queryParams := map[string]string{
		"team":   strconv.Itoa(teamId),
		"season": strconv.Itoa(Season),
	}

	apiUrl, err := url.Parse(fmt.Sprintf("%s/v3/fixtures", h.rapidApi.Config.Football.BaseUrl))
	if err != nil {
		return FixturesResponse{}, fmt.Errorf("faiiled to parse Football Api games url: %w", err)
	}

	response, err := h.rapidApi.BaseRequest().SetQueryParams(queryParams).Get(apiUrl.String())
	if err != nil {
		return FixturesResponse{}, fmt.Errorf("request failed to retrieve Football games: %w", err)
	}

	h.logger.Debug("Fotball Api response", "response", response)

	return h.parseGamesResponse(response.Body())
}

func (h *Handler) parseGamesResponse(input []byte) (FixturesResponse, error) {
	var data FixturesResponse
	if err := json.Unmarshal(input, &data); err != nil {
		return data, fmt.Errorf("failed to unmarshall Football games: %w", err)
	}

	return data, nil
}

func (h *Handler) createCalendar(games []Fixture) calendar.Calendar {
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

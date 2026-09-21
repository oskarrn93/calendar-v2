package nba

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
)

type GameDate struct {
	Start time.Time `json:"start"`
}

type Team struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GameTeams struct {
	Visitors Team `json:"visitors"`
	Home     Team `json:"home"`
}

type Game struct {
	Id    int       `json:"id"`
	Date  GameDate  `json:"date"`
	Teams GameTeams `json:"teams"`
}

type GamesResponse struct {
	Get      string `json:"get"`
	Results  int    `json:"results"`
	Response []Game `json:"response"`
}

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

	calendar := h.createCalendar(games)
	calendarData, err := calendar.Export()
	if err != nil {
		return fmt.Errorf("failed to export NBA calendar: %w", err)
	}

	if err := h.storage.Upload(ctx, "nba.ics", calendarData); err != nil {
		return fmt.Errorf("failed to upload NBA file: %w", err)
	}

	return nil
}

func (h *Handler) getGamesByTeam(teamId TeamID) (GamesResponse, error) {
	// TODO: Add support for multiple teams
	queryParams := map[string]string{
		"team":   strconv.Itoa(int(teamId)),
		"season": strconv.Itoa(Season),
	}

	apiUrl := fmt.Sprintf("%s/games", h.rapidApi.Config.NBA.BaseUrl)

	response, err := h.rapidApi.BaseRequest().SetQueryParams(queryParams).Get(apiUrl)
	if err != nil {
		return GamesResponse{}, fmt.Errorf("request failed to retrieve NBA games: %w", err)
	}

	h.logger.Debug("NBA Api response", "response", response)

	var data GamesResponse
	if err := json.Unmarshal(response.Body(), &data); err != nil {
		return GamesResponse{}, fmt.Errorf("failed to unmarshall NBA games: %w", err)
	}

	return data, nil
}

func (h *Handler) GetGames(teamIds []TeamID) ([]Game, error) {
	games := []Game{}

	for _, teamId := range teamIds {
		data, err := h.getGamesByTeam(teamId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve NBA games for team id %d: %w", teamId, err)
		}
		games = append(games, data.Response...)
	}

	return games, nil
}

func (h *Handler) createCalendar(games []Game) *calendar.Calendar {
	cal := calendar.New("NBA")

	for _, game := range games {
		cal.AddEvent(calendar.Event{
			Id:        fmt.Sprintf("nba-%d", game.Id),
			Title:     fmt.Sprintf("%s - %s", game.Teams.Home.Name, game.Teams.Visitors.Name),
			StartDate: game.Date.Start,
			EndDate:   game.Date.Start.Add(2 * time.Hour),
		})
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

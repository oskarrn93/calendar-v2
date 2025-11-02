package esport

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
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
	events, err := h.GetEvents([]SportID{EsportSportID})
	if err != nil {
		return err
	}

	h.logger.Debug("Esport events", "events", events)

	calendar := h.createCalendar(events)
	calendarData, err := calendar.Export()
	if err != nil {
		return fmt.Errorf("failed to export Esport calendar: %w", err)
	}

	if err := h.storage.Upload(ctx, "esport.ics", calendarData, h.logger); err != nil {
		return fmt.Errorf("failed to upload Esport file: %w", err)
	}

	return nil
}

type Event struct {
	ID         int    `json:"event_id"`
	LeagueName string `json:"league_name"`
	Starts     string `json:"starts"`
	Last       int64  `json:"last"`
	Home       string `json:"home"`
	Away       string `json:"away"`
}

func (e *Event) IsCS2() bool {
	return strings.Contains(strings.ToLower(e.LeagueName), "cs2")
}

func (e *Event) HasTeam(team string) bool {
	team = strings.ToLower(team)
	return strings.Contains(strings.ToLower(e.Home), team) || strings.Contains(strings.ToLower(e.Away), team)
}

func (e *Event) HasTeams(teams []string) bool {
	for _, team := range teams {
		if e.HasTeam(team) {
			return true
		}
	}
	return false
}

type EventsResponse struct {
	Events []Event `json:"events"`
}

func (h *Handler) GetEvents(sportIDs []SportID) ([]Event, error) {
	events := []Event{}

	for _, sportID := range sportIDs {
		response, err := h.getEventsBySport(int(sportID))
		if err != nil {
			return nil, err
		}

		h.logger.Debug("Retrieved Esport games", "sportID", sportID, "events", response.Events)
		events = append(events, response.Events...)
	}

	return events, nil
}

func (h *Handler) getEventsBySport(sportID int) (EventsResponse, error) {
	// API docs https://rapidapi.com/tipsters/api/pinnacle-odds

	/*
		curl --request GET
		--url 'https://pinnacle-odds.p.rapidapi.com/kit/v1/markets?sport_id=10'
		--header 'x-rapidapi-key: REPLACE_ME'
	*/

	queryParams := map[string]string{
		"sport_id": strconv.Itoa(sportID),
	}

	apiUrl, err := url.Parse(fmt.Sprintf("%s/kit/v1/markets", h.rapidApi.Config.Esport.BaseUrl))
	if err != nil {
		return EventsResponse{}, fmt.Errorf("faiiled to parse Esport Api games url: %w", err)
	}

	response, err := h.rapidApi.BaseRequest().SetQueryParams(queryParams).Get(apiUrl.String())
	if err != nil {
		return EventsResponse{}, fmt.Errorf("request failed to retrieve Esport games: %w", err)
	}

	h.logger.Debug("Fotball Api response", "response", response)

	return h.parseGamesResponse(response.Body())
}

func (h *Handler) parseGamesResponse(input []byte) (EventsResponse, error) {
	var data EventsResponse
	if err := json.Unmarshal(input, &data); err != nil {
		return data, fmt.Errorf("failed to unmarshall Esport games: %w", err)
	}

	return data, nil
}

func (h *Handler) createCalendar(events []Event) calendar.Calendar {
	cal := calendar.New("Esport")
	for _, event := range events {
		if !event.IsCS2() {
			continue
		}

		if !event.HasTeams(TeamsOfInterest) {
			continue
		}

		startDate, err := time.Parse("2006-01-02T15:04:05", event.Starts)
		if err != nil {
			h.logger.Warn("Failed to parse Esport event start date", "eventID", event.ID, "starts", event.Starts, "error", err)
			continue
		}

		newEvent := calendar.Event{
			Id:        fmt.Sprintf("esport-%d", event.ID),
			Title:     fmt.Sprintf("%s - %s", event.Home, event.Away),
			StartDate: startDate,
			EndDate:   time.Unix(event.Last, 0),
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

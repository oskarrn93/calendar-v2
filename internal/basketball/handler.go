package basketball

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
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

	h.logger.Debug("Basketball games", "games", games)

	calendar := h.createCalendar(games)
	calendarData, err := calendar.Export()
	if err != nil {
		return fmt.Errorf("failed to export Basketball calendar: %w", err)
	}

	if err := h.storage.Upload(ctx, "basketball.ics", calendarData, h.logger); err != nil {
		return fmt.Errorf("failed to upload Basketball file: %w", err)
	}

	return nil
}

type EventTeam struct {
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Gender    string `json:"gender"`
}

func (e EventTeam) PrettyName() string {
	if e.ShortName != "" {
		return e.ShortName
	}

	return e.Name
}

type EventTournament struct {
	Name string `json:"name"`
}

type Event struct {
	ID             uint64          `json:"id"`
	Tournament     EventTournament `json:"tournament"`
	HomeTeam       EventTeam       `json:"homeTeam"`
	AwayTeam       EventTeam       `json:"awayTeam"`
	StartTimestamp uint64          `json:"startTimestamp"`
}

type EventsResponse struct {
	Events      []Event `json:"events"`
	HasNextPage bool    `json:"hasNextPage"`
}

func (h *Handler) GetGames(teamIDs []TeamID) ([]Event, error) {
	events := []Event{}

	for _, teamID := range teamIDs {
		response, err := h.getEventsByTeam(int(teamID))
		if err != nil {
			return nil, err
		}

		h.logger.Debug("Retrieved Basketball games", "teamId", teamID, "events", response.Events)
		events = append(events, response.Events...)
	}

	return events, nil
}

func (h *Handler) getEventsByTeam(teamId int) (EventsResponse, error) {
	/*
		curl --request GET
		--url https://sportapi7.p.rapidapi.com/api/v1/team/3540/events/next/1
		--header 'x-rapidapi-key: REPLACE_ME'
	*/

	apiUrl, err := url.Parse(fmt.Sprintf("%s/api/v1/team/%d/events/next/1", h.rapidApi.Config.Basketball.BaseUrl, teamId))
	if err != nil {
		return EventsResponse{}, fmt.Errorf("faiiled to parse Basketball Api games url: %w", err)
	}

	response, err := h.rapidApi.BaseRequest().Get(apiUrl.String())
	if err != nil {
		return EventsResponse{}, fmt.Errorf("request failed to retrieve Basketball games: %w", err)
	}

	h.logger.Debug("Basketball Api response", "response", response)

	return h.parseEventsResponse(response.Body())
}

func (h *Handler) parseEventsResponse(input []byte) (EventsResponse, error) {
	var data EventsResponse
	if err := json.Unmarshal(input, &data); err != nil {
		return data, fmt.Errorf("failed to unmarshall Basketball games: %w", err)
	}

	return data, nil
}

func (h *Handler) createCalendar(events []Event) calendar.Calendar {
	cal := calendar.New("Basketball")
	for _, event := range events {

		startTime := time.Unix(int64(event.StartTimestamp), 0)

		newEvent := calendar.Event{
			Id:        fmt.Sprintf("basketball-%d", event.ID),
			Title:     fmt.Sprintf("%s: %s - %s", event.Tournament.Name, event.HomeTeam.PrettyName(), event.AwayTeam.PrettyName()),
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

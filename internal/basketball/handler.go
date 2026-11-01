package basketball

import (
	"context"
	"fmt"
	"time"

	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/httpclient"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
)

type Handler struct {
	rapidApi rapidapi.RapidApi
	storage  awsutil.Storage
}

func (h *Handler) Handler(ctx context.Context) error {
	games, err := h.GetGames(ctx, TeamIDs)
	if err != nil {
		return err
	}

	return createCalendar(games).Upload(ctx, h.storage, "basketball.ics")
}

type EventTeam struct {
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
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
	ID             int64           `json:"id"`
	Tournament     EventTournament `json:"tournament"`
	HomeTeam       EventTeam       `json:"homeTeam"`
	AwayTeam       EventTeam       `json:"awayTeam"`
	StartTimestamp int64           `json:"startTimestamp"`
}

type EventsResponse struct {
	Events      []Event `json:"events"`
	HasNextPage bool    `json:"hasNextPage"`
}

func (h *Handler) GetGames(ctx context.Context, teamIDs []TeamID) ([]Event, error) {
	events := []Event{}

	for _, teamID := range teamIDs {
		for page := 0; page < maxPages; page++ {
			response, err := h.getEventsByTeam(ctx, int(teamID), page)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve Basketball games for team id %d page %d: %w", teamID, page, err)
			}

			events = append(events, response.Events...)

			if !response.HasNextPage {
				break
			}
		}
	}

	return events, nil
}

func (h *Handler) getEventsByTeam(ctx context.Context, teamId int, page int) (EventsResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v1/team/%d/events/next/%d", BaseUrl, teamId, page)

	return httpclient.GetJSON[EventsResponse](h.rapidApi.BaseRequest(ctx), apiUrl)
}

func createCalendar(events []Event) *calendar.Calendar {
	cal := calendar.New("Basketball")
	for _, event := range events {

		startTime := time.Unix(event.StartTimestamp, 0)

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

func NewHandler(rapidApi rapidapi.RapidApi, storage awsutil.Storage) *Handler {
	return &Handler{
		rapidApi: rapidApi,
		storage:  storage,
	}
}

package football

import (
	"context"
	"fmt"
	"strconv"
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

	return createCalendar(games).Upload(ctx, h.storage, "football.ics")
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
	Response []Fixture `json:"response"`
}

func (h *Handler) GetGames(ctx context.Context, teamIds []TeamID) ([]Fixture, error) {
	games := []Fixture{}

	for _, teamId := range teamIds {
		teamGames, err := h.getGamesByTeam(ctx, teamId)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve Football games for team id %d: %w", teamId, err)
		}

		games = append(games, teamGames...)
	}

	return games, nil
}

func (h *Handler) getGamesByTeam(ctx context.Context, teamId TeamID) ([]Fixture, error) {
	queryParams := map[string]string{
		"team":   strconv.Itoa(int(teamId)),
		"season": strconv.Itoa(Season),
	}

	data, err := httpclient.GetJSON[FixturesResponse](h.rapidApi.BaseRequest(ctx).SetQueryParams(queryParams), BaseUrl+"/v3/fixtures")

	return data.Response, err
}

func createCalendar(games []Fixture) *calendar.Calendar {
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

func NewHandler(rapidApi rapidapi.RapidApi, storage awsutil.Storage) *Handler {
	return &Handler{
		rapidApi: rapidApi,
		storage:  storage,
	}
}

package football

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/httpclient"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"golang.org/x/sync/errgroup"
)

// Keeps concurrent requests under RapidAPI's per-second rate limit.
const maxConcurrentRequests = 3

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
	Errors   json.RawMessage `json:"errors"`
	Response []Fixture       `json:"response"`
}

// API-Sports reports failures such as an exhausted quota or a bad key with a
// 200 status and an empty response. errors is [] on success and an object
// otherwise.
func (r FixturesResponse) err() error {
	switch e := string(bytes.TrimSpace(r.Errors)); e {
	case "", "null", "[]", "{}":
		return nil
	default:
		return fmt.Errorf("api error: %s", e)
	}
}

func (h *Handler) GetGames(ctx context.Context, teamIds []TeamID) ([]Fixture, error) {
	teamGames := make([][]Fixture, len(teamIds))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrentRequests)

	for i, teamId := range teamIds {
		g.Go(func() error {
			games, err := h.getGamesByTeam(ctx, teamId)
			if err != nil {
				return fmt.Errorf("failed to retrieve Football games for team id %d: %w", teamId, err)
			}
			teamGames[i] = games
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return slices.Concat(teamGames...), nil
}

func (h *Handler) getGamesByTeam(ctx context.Context, teamId TeamID) ([]Fixture, error) {
	queryParams := map[string]string{
		"team":   strconv.Itoa(int(teamId)),
		"season": strconv.Itoa(Season(teamId, time.Now())),
	}

	data, err := httpclient.GetJSON[FixturesResponse](h.rapidApi.BaseRequest(ctx).SetQueryParams(queryParams), BaseUrl+"/v3/fixtures")
	if err != nil {
		return nil, err
	}

	return data.Response, data.err()
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

package nba

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/httpclient"
	"github.com/oskarrn93/calendar-v2/internal/rapidapi"
	"golang.org/x/sync/errgroup"
)

// Keeps concurrent requests under RapidAPI's per-second rate limit.
const maxConcurrentRequests = 3

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

type GameStatus struct {
	Long string `json:"long"`
}

type Game struct {
	Id     int        `json:"id"`
	Date   GameDate   `json:"date"`
	Status GameStatus `json:"status"`
	Teams  GameTeams  `json:"teams"`
}

func (g Game) isCancelled() bool {
	switch strings.ToLower(g.Status.Long) {
	case "postponed", "canceled", "cancelled":
		return true
	default:
		return false
	}
}

type GamesResponse struct {
	Errors   json.RawMessage `json:"errors"`
	Response []Game          `json:"response"`
}

// API-Sports reports failures such as an exhausted quota or a bad key with a
// 200 status and an empty response. errors is [] on success and an object
// otherwise.
func (r GamesResponse) err() error {
	switch e := string(bytes.TrimSpace(r.Errors)); e {
	case "", "null", "[]", "{}":
		return nil
	default:
		return fmt.Errorf("api error: %s", e)
	}
}

type Handler struct {
	rapidApi rapidapi.RapidApi
	storage  awsutil.Storage
}

func (h *Handler) Handler(ctx context.Context) error {
	games, err := h.GetGames(ctx, TeamIDs)
	if err != nil {
		return err
	}

	return createCalendar(games).Upload(ctx, h.storage, "nba.ics")
}

func (h *Handler) getGamesByTeam(ctx context.Context, teamId TeamID) ([]Game, error) {
	queryParams := map[string]string{
		"team":   strconv.Itoa(int(teamId)),
		"season": strconv.Itoa(Season(time.Now())),
	}

	data, err := httpclient.GetJSON[GamesResponse](h.rapidApi.BaseRequest(ctx).SetQueryParams(queryParams), BaseUrl+"/games")
	if err != nil {
		return nil, err
	}

	return data.Response, data.err()
}

func (h *Handler) GetGames(ctx context.Context, teamIds []TeamID) ([]Game, error) {
	teamGames := make([][]Game, len(teamIds))

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrentRequests)

	for i, teamId := range teamIds {
		g.Go(func() error {
			games, err := h.getGamesByTeam(ctx, teamId)
			if err != nil {
				return fmt.Errorf("failed to retrieve NBA games for team id %d: %w", teamId, err)
			}
			teamGames[i] = games
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// A team can come back empty with no error. Fail so the previous file is
	// kept, unless every team is empty, which is just the off-season.
	var emptyTeams []TeamID
	for i, games := range teamGames {
		if len(games) == 0 {
			emptyTeams = append(emptyTeams, teamIds[i])
		}
	}
	if len(emptyTeams) > 0 && len(emptyTeams) < len(teamIds) {
		return nil, fmt.Errorf("no NBA games returned for team ids %v while other teams have games", emptyTeams)
	}

	return slices.Concat(teamGames...), nil
}

func createCalendar(games []Game) *calendar.Calendar {
	cal := calendar.New("NBA")

	for _, game := range games {
		if game.isCancelled() {
			continue
		}

		cal.AddEvent(calendar.Event{
			Id:        fmt.Sprintf("nba-%d", game.Id),
			Title:     fmt.Sprintf("%s - %s", game.Teams.Home.Name, game.Teams.Visitors.Name),
			StartDate: game.Date.Start,
			EndDate:   game.Date.Start.Add(2 * time.Hour),
		})
	}

	return cal
}

func NewHandler(rapidApi rapidapi.RapidApi, storage awsutil.Storage) *Handler {
	return &Handler{
		rapidApi: rapidApi,
		storage:  storage,
	}
}

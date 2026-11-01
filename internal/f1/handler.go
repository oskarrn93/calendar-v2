package f1

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/oskarrn93/calendar-v2/internal/awsutil"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/httpclient"
)

const (
	BaseUrl    = "https://api.jolpi.ca"
	websiteUrl = "https://www.formula1.com"
)

type Session struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Race struct {
	Season           string   `json:"season"`
	Round            string   `json:"round"`
	RaceName         string   `json:"raceName"`
	Date             string   `json:"date"`
	Time             string   `json:"time"`
	Qualifying       *Session `json:"Qualifying"`
	Sprint           *Session `json:"Sprint"`
	SprintQualifying *Session `json:"SprintQualifying"`
}

type racesResponse struct {
	MRData struct {
		RaceTable struct {
			Races []Race `json:"Races"`
		} `json:"RaceTable"`
	} `json:"MRData"`
}

type Handler struct {
	httpClient *resty.Client
	storage    awsutil.Storage
	logger     *slog.Logger
}

func (h *Handler) Handler(ctx context.Context) error {
	// Next season's schedule is usually published in the autumn; until then the
	// API returns it as an empty list.
	year := time.Now().UTC().Year()
	races, err := h.GetRaces(ctx, []int{year, year + 1})
	if err != nil {
		return err
	}

	return h.createCalendar(races).Upload(ctx, h.storage, "f1.ics")
}

func (h *Handler) GetRaces(ctx context.Context, seasons []int) ([]Race, error) {
	races := []Race{}

	for _, season := range seasons {
		seasonRaces, err := h.getRacesBySeason(ctx, season)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve F1 races for season %d: %w", season, err)
		}
		races = append(races, seasonRaces...)
	}

	return races, nil
}

// API docs https://github.com/jolpica/jolpica-f1/blob/main/docs/README.md
func (h *Handler) getRacesBySeason(ctx context.Context, season int) ([]Race, error) {
	apiUrl := fmt.Sprintf("%s/ergast/f1/%d/races/", BaseUrl, season)

	data, err := httpclient.GetJSON[racesResponse](h.httpClient.R().SetContext(ctx).SetQueryParam("limit", "100"), apiUrl)
	if err != nil {
		return nil, err
	}

	return data.MRData.RaceTable.Races, nil
}

func (h *Handler) createCalendar(races []Race) *calendar.Calendar {
	cal := calendar.New("F1")

	for _, race := range races {
		sessions := []struct {
			key      string
			name     string
			session  *Session
			duration time.Duration
		}{
			{"sprint-qualifying", "Sprint Qualifying", race.SprintQualifying, time.Hour},
			{"sprint", "Sprint", race.Sprint, time.Hour},
			{"qualifying", "Qualifying", race.Qualifying, time.Hour},
			{"race", "Race", &Session{Date: race.Date, Time: race.Time}, 2 * time.Hour},
		}

		for _, s := range sessions {
			if s.session == nil {
				continue
			}

			start, err := time.Parse(time.RFC3339, s.session.Date+"T"+s.session.Time)
			if err != nil {
				h.logger.Warn("Failed to parse F1 session start", "race", race.RaceName, "session", s.name, "error", err)
				continue
			}

			cal.AddEvent(calendar.Event{
				Id:        fmt.Sprintf("f1-%s-%s-%s", race.Season, race.Round, s.key),
				Title:     fmt.Sprintf("F1 %s - %s", race.RaceName, s.name),
				StartDate: start,
				EndDate:   start.Add(s.duration),
				Url:       websiteUrl,
			})
		}
	}

	return cal
}

func NewHandler(httpClient *resty.Client, storage awsutil.Storage, logger *slog.Logger) *Handler {
	return &Handler{
		httpClient: httpClient,
		storage:    storage,
		logger:     logger,
	}
}

package main

import (
	"context"
	"fmt"
	"time"
)

func NbaHandler(ctx context.Context, nbaApi NBAApi, storage S3Storage) error {
	calendar := NewCalendar("NBA")

	games := nbaApi.getGames()
	for _, game := range games.Response {
		newEvent := CalendarEvent{
			Id:        fmt.Sprintf("nba-%d", game.Id),
			Title:     fmt.Sprintf("%s - %s", game.Teams.Home.Name, game.Teams.Visitors.Name),
			StartDate: game.Date.Start,
			EndDate:   game.Date.Start.Add(2 * time.Hour),
		}
		calendar.AddEvent(newEvent)
	}

	if err := storage.upload(ctx, "nba.ics", calendar.Export()); err != nil {
		return fmt.Errorf("failed to upload NBA file: %w", err)
	}

	return nil
}

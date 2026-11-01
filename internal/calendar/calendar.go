package calendar

import (
	"bytes"
	"context"
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/oskarrn93/calendar-v2/internal/awsutil"
)

type Event struct {
	Id        string
	Title     string
	StartDate time.Time
	EndDate   time.Time
	Url       string
}

type Calendar struct {
	calendar  *ics.Calendar
	createdAt time.Time
	eventIDs  map[string]struct{}
}

func New(name string) *Calendar {
	calendar := ics.NewCalendar()
	calendar.SetProductId(fmt.Sprintf("-//%s", name))
	calendar.SetName(name)

	return &Calendar{
		calendar:  calendar,
		createdAt: time.Now(),
		eventIDs:  map[string]struct{}{},
	}
}

// AddEvent ignores events whose Id was already added: when two followed teams
// play each other the same game is returned once per team, and duplicate UIDs
// make calendar clients drop or mangle the event.
func (cal *Calendar) AddEvent(newEvent Event) {
	if _, ok := cal.eventIDs[newEvent.Id]; ok {
		return
	}
	cal.eventIDs[newEvent.Id] = struct{}{}

	icsEvent := cal.calendar.AddEvent(newEvent.Id)

	icsEvent.SetDtStampTime(cal.createdAt)
	icsEvent.SetSummary(newEvent.Title)
	icsEvent.SetStartAt(newEvent.StartDate)
	icsEvent.SetEndAt(newEvent.EndDate)

	if newEvent.Url != "" {
		icsEvent.SetURL(newEvent.Url)
		// Google Calendar ignores the URL property, so repeat it where every client shows it.
		icsEvent.SetDescription(newEvent.Url)
	}
}

func (cal *Calendar) Export() ([]byte, error) {
	var data bytes.Buffer
	if err := cal.calendar.SerializeTo(&data); err != nil {
		return nil, fmt.Errorf("failed to serialize calendar: %w", err)
	}
	return data.Bytes(), nil
}

func (cal *Calendar) Upload(ctx context.Context, storage awsutil.Storage, filename string) error {
	data, err := cal.Export()
	if err != nil {
		return fmt.Errorf("failed to export %s: %w", filename, err)
	}

	if err := storage.Upload(ctx, filename, data); err != nil {
		return fmt.Errorf("failed to upload %s: %w", filename, err)
	}

	return nil
}

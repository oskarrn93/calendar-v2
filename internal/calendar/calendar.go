package calendar

import (
	"bytes"
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Event struct {
	Id        string
	Title     string
	StartDate time.Time
	EndDate   time.Time
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
}

func (cal *Calendar) Export() ([]byte, error) {
	var data bytes.Buffer
	if err := cal.calendar.SerializeTo(&data); err != nil {
		return nil, fmt.Errorf("failed to serialize calendar: %w", err)
	}
	return data.Bytes(), nil
}

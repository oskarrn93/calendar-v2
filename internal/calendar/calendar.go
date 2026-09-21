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
	calendar *ics.Calendar
}

func New(name string) *Calendar {
	calendar := ics.NewCalendar()
	calendar.SetProductId(fmt.Sprintf("-//%s", name))
	calendar.SetName(name)

	return &Calendar{calendar: calendar}
}

func (cal *Calendar) AddEvent(newEvent Event) {
	icsEvent := cal.calendar.AddEvent(newEvent.Id)

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

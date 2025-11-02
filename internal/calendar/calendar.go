package calendar

import (
	"bytes"
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Calendar interface {
	AddEvent(event Event)
	Export() ([]byte, error)
}

type Event struct {
	Id        string
	Title     string
	StartDate time.Time
	EndDate   time.Time
}

type ICSCalendar struct {
	calendar *ics.Calendar
}

func (cal ICSCalendar) AddEvent(newEvent Event) {
	icsEvent := cal.calendar.AddEvent(newEvent.Id)

	icsEvent.SetSummary(newEvent.Title)
	icsEvent.SetStartAt(newEvent.StartDate)
	icsEvent.SetEndAt(newEvent.EndDate)
}

func (cal ICSCalendar) Export() ([]byte, error) {
	var data bytes.Buffer
	err := cal.calendar.SerializeTo(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialze calendar: %w", err)
	}
	return data.Bytes(), nil
}

func New(name string) Calendar {
	calendar := *ics.NewCalendar()
	calendar.SetProductId(fmt.Sprintf("-//%s", name))
	calendar.SetName(name)

	return &ICSCalendar{
		calendar: &calendar,
	}
}

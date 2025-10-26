package main

import (
	"bytes"
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Calendar interface {
	AddEvent(event CalendarEvent)
	Export() []byte
}

type CalendarEvent struct {
	Id        string
	Title     string
	StartDate time.Time
	EndDate   time.Time
}

type ICSCalendar struct {
	calendar *ics.Calendar
}

func (cal ICSCalendar) AddEvent(newEvent CalendarEvent) {
	icsEvent := cal.calendar.AddEvent(newEvent.Id)

	icsEvent.SetSummary(newEvent.Title)
	icsEvent.SetStartAt(newEvent.StartDate)
	icsEvent.SetEndAt(newEvent.EndDate)
}

func (cal ICSCalendar) Export() []byte {
	var data bytes.Buffer
	cal.calendar.SerializeTo(&data)
	return data.Bytes()
}

func NewCalendar(name string) Calendar {
	calendar := *ics.NewCalendar()
	calendar.SetProductId(fmt.Sprintf("-//%s", name))
	calendar.SetName(name)

	return &ICSCalendar{
		calendar: &calendar,
	}
}

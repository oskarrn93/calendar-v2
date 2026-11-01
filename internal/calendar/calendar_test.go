package calendar_test

import (
	"bytes"
	"testing"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/oskarrn93/calendar-v2/internal/calendar"
	"github.com/oskarrn93/calendar-v2/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestAddEventSkipsDuplicateIDs(t *testing.T) {
	cal := calendar.New("Test")
	start := time.Date(2026, 1, 1, 18, 0, 0, 0, time.UTC)

	cal.AddEvent(calendar.Event{Id: "game-1", Title: "A - B", StartDate: start, EndDate: start.Add(2 * time.Hour)})
	cal.AddEvent(calendar.Event{Id: "game-1", Title: "A - B", StartDate: start, EndDate: start.Add(2 * time.Hour)})
	cal.AddEvent(calendar.Event{Id: "game-2", Title: "C - D", StartDate: start, EndDate: start.Add(2 * time.Hour)})

	data, err := cal.Export()
	require.NoError(t, err)

	parsed, err := ics.ParseCalendar(bytes.NewReader(data))
	require.NoError(t, err)

	events := parsed.Events()
	require.Len(t, events, 2)

	for _, event := range events {
		_, err := event.GetDtStampTime()
		require.NoError(t, err, "event %s is missing DTSTAMP", event.Id())
	}
}

func TestAddEventSetsUrl(t *testing.T) {
	cal := calendar.New("Test")
	start := time.Date(2026, 1, 1, 18, 0, 0, 0, time.UTC)

	cal.AddEvent(calendar.Event{Id: "with-url", Title: "A", StartDate: start, EndDate: start.Add(time.Hour), Url: "https://example.com/race"})
	cal.AddEvent(calendar.Event{Id: "without-url", Title: "B", StartDate: start, EndDate: start.Add(time.Hour)})

	data, err := cal.Export()
	require.NoError(t, err)

	parsed, err := ics.ParseCalendar(bytes.NewReader(data))
	require.NoError(t, err)

	events := parsed.Events()
	require.Len(t, events, 2)

	withUrl := events[0].GetProperty(ics.ComponentPropertyUrl)
	require.NotNil(t, withUrl)
	require.Equal(t, "https://example.com/race", withUrl.Value)
	require.Equal(t, "https://example.com/race", events[0].GetProperty(ics.ComponentPropertyDescription).Value)

	require.Nil(t, events[1].GetProperty(ics.ComponentPropertyUrl))
	require.Nil(t, events[1].GetProperty(ics.ComponentPropertyDescription))
}

func TestUploadRefusesEmptyCalendar(t *testing.T) {
	storage := &testutil.MemoryStorage{}

	require.ErrorIs(t, calendar.New("Test").Upload(t.Context(), storage, "test.ics"), calendar.ErrNoEvents)
	require.Empty(t, storage.Files)
}

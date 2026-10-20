package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestReadEventsByPeriod(t *testing.T) {
	store := New()
	ctx := context.Background()

	day := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	events := []storage.Event{
		{
			Title:     "today",
			DateStart: day.Add(10 * time.Hour),
			DateEnd:   day.Add(11 * time.Hour),
			UserId:    1,
		},
		{
			Title:     "tomorrow",
			DateStart: day.AddDate(0, 0, 1).Add(10 * time.Hour),
			DateEnd:   day.AddDate(0, 0, 1).Add(11 * time.Hour),
			UserId:    1,
		},
	}

	for _, event := range events {
		_, err := store.CreateEvent(event, ctx)
		require.NoError(t, err)
	}

	from, to, err := storage.PeriodRange(day, storage.PeriodDay)
	require.NoError(t, err)

	result, err := store.ReadEventsByPeriod(from, to, ctx)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "today", result[0].Title)

	from, to, err = storage.PeriodRange(day, storage.PeriodWeek)
	require.NoError(t, err)

	result, err = store.ReadEventsByPeriod(from, to, ctx)
	require.NoError(t, err)
	require.Len(t, result, 2)
}

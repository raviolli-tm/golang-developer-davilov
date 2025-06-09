package memorystorage

import (
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	calendar := New()

	tests := []struct {
		operation int
		input     storage.Event
		expected  []storage.Event
		err       error
	}{
		{operation: 0, input: storage.Event{
			Title:           "test title",
			DateStart:       time.Now(),
			DateEnd:         time.Now().Add(15 * time.Minute),
			Description:     "test description",
			UserId:          1,
			EventNotifyTime: 15,
		}},
		{operation: 1, expected: []storage.Event{{
			Title:           "test title",
			DateStart:       time.Now(),
			DateEnd:         time.Now().Add(15 * time.Minute),
			Description:     "test description",
			UserId:          1,
			EventNotifyTime: 15,
		}}},
		{operation: 3, input: storage.Event{
			Title:           "test title(updated)",
			DateStart:       time.Now(),
			DateEnd:         time.Now().Add(15 * time.Minute),
			Description:     "test description(updated)",
			UserId:          1,
			EventNotifyTime: 15,
		}},
		{operation: 1, expected: []storage.Event{{
			Title:           "test title(updated)",
			DateStart:       time.Now(),
			DateEnd:         time.Now().Add(15 * time.Minute),
			Description:     "test description(updated)",
			UserId:          1,
			EventNotifyTime: 15,
		}}},
		{operation: 2, err: appErrors.ErrIdDoesNotExist},
		{operation: 1, expected: []storage.Event{}},
		{operation: 0, input: storage.Event{
			Title:           "test title",
			DateStart:       time.Now(),
			DateEnd:         time.Now().Add(15 * time.Minute),
			Description:     "test description",
			UserId:          1,
			EventNotifyTime: 15,
		}},
		{operation: 0, input: storage.Event{
			Title:           "test title 2",
			DateStart:       time.Now(),
			DateEnd:         time.Now().Add(15 * time.Minute),
			Description:     "test description 2",
			UserId:          1,
			EventNotifyTime: 15,
		}},
		{operation: 1, expected: []storage.Event{{
			Title:           "test title",
			DateStart:       time.Now(),
			DateEnd:         time.Now().Add(15 * time.Minute),
			Description:     "test description",
			UserId:          1,
			EventNotifyTime: 15,
		}, {
			Title:           "test title 2",
			DateStart:       time.Now(),
			DateEnd:         time.Now().Add(15 * time.Minute),
			Description:     "test description 2",
			UserId:          1,
			EventNotifyTime: 15,
		}}},
	}

	var id uuid.UUID
	for testId, ts := range tests {
		switch ts.operation {
		case 0:
			{
				err := calendar.CreateEvent(ts.input)
				require.NoError(t, err)

			}
		case 2:
			{
				err := calendar.DeleteEvent(id)
				require.NoError(t, err)
			}
		case 3:
			{
				err := calendar.UpdateEvent(id, ts.input)
				require.NoError(t, err)
			}
		case 1:
			{
				events, err := calendar.ReadEvents()
				require.NoError(t, err)
				require.Equal(t, len(ts.expected), len(events))
				if len(events) > 0 {
					id = events[0].ID
				}
				for _, eventExpected := range ts.expected {
					for _, event := range events {
						if eventExpected.Title == event.Title {
							require.Equal(t, eventExpected.Title, event.Title)
							require.Equal(t, eventExpected.Description, event.Description)
							require.Equal(t, eventExpected.UserId, event.UserId)
						}
					}
				}

			}

		}
		fmt.Printf("Test %d passed\n", testId)

	}
}

func TestStorageErr(t *testing.T) {
	calendar := New()

	tests := []struct {
		operation int
		err       error
	}{
		{operation: 2, err: appErrors.ErrIdDoesNotExist},
		{operation: 3, err: appErrors.ErrIdDoesNotExist},
	}

	for testId, ts := range tests {
		switch ts.operation {
		case 0:
			{
			}
		case 1:
			{
			}
		case 2:
			{
				err := calendar.DeleteEvent(uuid.New())
				require.ErrorIs(t, err, ts.err)
			}

		case 3:
			{
				err := calendar.UpdateEvent(uuid.New(), storage.Event{
					Title:           "test title",
					DateStart:       time.Now(),
					DateEnd:         time.Now().Add(15 * time.Minute),
					Description:     "test description",
					UserId:          1,
					EventNotifyTime: 15,
				})
				require.ErrorIs(t, err, ts.err)
			}

		}
		fmt.Printf("Test %d passed\n", testId)
	}
}

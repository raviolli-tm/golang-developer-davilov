package memorystorage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func validEvent() storage.Event {
	return storage.Event{
		Title:           "test title",
		DateStart:       time.Now(),
		DateEnd:         time.Now().Add(15 * time.Minute),
		Description:     "test description",
		UserId:          1,
		EventNotifyTime: 15,
	}
}

func TestStorageCRUD(t *testing.T) {
	store := New()
	ctx := context.Background()

	created, err := store.CreateEvent(validEvent(), ctx)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, created.ID)

	events, err := store.ReadEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, created.ID, events[0].ID)

	updated, err := store.UpdateEvent(created.ID, storage.Event{
		Title:           "updated title",
		DateStart:       time.Now(),
		DateEnd:         time.Now().Add(time.Hour),
		Description:     "updated description",
		UserId:          2,
		EventNotifyTime: 5,
	}, ctx)
	require.NoError(t, err)
	require.Equal(t, "updated title", updated.Title)

	events, err = store.ReadEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "updated title", events[0].Title)

	deletedID, err := store.DeleteEvent(created.ID, ctx)
	require.NoError(t, err)
	require.Equal(t, created.ID, deletedID)

	events, err = store.ReadEvents(ctx)
	require.NoError(t, err)
	require.Empty(t, events)
}

func TestStorageErrors(t *testing.T) {
	store := New()
	ctx := context.Background()

	t.Run("update missing event", func(t *testing.T) {
		_, err := store.UpdateEvent(uuid.New(), validEvent(), ctx)
		require.ErrorIs(t, err, appErrors.ErrIdDoesNotExist)
	})

	t.Run("delete missing event", func(t *testing.T) {
		_, err := store.DeleteEvent(uuid.New(), ctx)
		require.ErrorIs(t, err, appErrors.ErrIdDoesNotExist)
	})

	t.Run("validation error", func(t *testing.T) {
		_, err := store.CreateEvent(storage.Event{}, ctx)
		require.True(t, errors.Is(err, appErrors.ErrValidation))
	})
}

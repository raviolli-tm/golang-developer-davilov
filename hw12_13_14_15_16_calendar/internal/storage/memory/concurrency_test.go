package memorystorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorageConcurrency(t *testing.T) {
	store := New()
	ctx := context.Background()

	const (
		workers    = 10
		iterations = 100
	)

	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				event := storage.Event{
					Title:     "concurrent",
					DateStart: time.Now().Add(time.Hour),
					DateEnd:   time.Now().Add(2 * time.Hour),
					UserId:    1,
				}

				if _, err := store.CreateEvent(event, ctx); err != nil {
					errCh <- err
					return
				}
				if _, err := store.ReadEvents(ctx); err != nil {
					errCh <- err
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err)
	}

	events, err := store.ReadEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, workers*iterations)
}

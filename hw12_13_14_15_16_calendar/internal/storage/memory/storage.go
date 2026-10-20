package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	mu   sync.RWMutex
	data map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{data: make(map[uuid.UUID]storage.Event)}
}

func (s *Storage) UpdateEvent(id uuid.UUID, e storage.Event, _ context.Context) (storage.Event, error) {
	if err := e.Validate(); err != nil {
		return storage.Event{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return storage.Event{}, appErrors.ErrIdDoesNotExist
	}

	e.ID = id
	s.data[id] = e

	return e, nil
}

func (s *Storage) ReadEvents(_ context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0, len(s.data))
	for _, event := range s.data {
		result = append(result, event)
	}

	return result, nil
}

// ReadEventsByPeriod returns events whose start time belongs to [from, to).
func (s *Storage) ReadEventsByPeriod(from, to time.Time, _ context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0)
	for _, event := range s.data {
		if !event.DateStart.Before(from) && event.DateStart.Before(to) {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) DeleteEvent(id uuid.UUID, _ context.Context) (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return uuid.Nil, appErrors.ErrIdDoesNotExist
	}

	delete(s.data, id)

	return id, nil
}

func (s *Storage) CreateEvent(e storage.Event, _ context.Context) (storage.Event, error) {
	if err := e.Validate(); err != nil {
		return storage.Event{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	e.ID = uuid.New()
	s.data[e.ID] = e

	return e, nil
}

// Close is a no-op for the in-memory storage, it exists to satisfy the
// app.Storage interface.
func (s *Storage) Close(_ context.Context) error {
	return nil
}

package memorystorage

import (
	"context"
	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"sync"
)

type Storage struct {
	mu   sync.RWMutex //nolint:unused
	data map[uuid.UUID]storage.Event
	idx  uuid.UUID
}

func New() *Storage {
	return &Storage{mu: sync.RWMutex{}, data: make(map[uuid.UUID]storage.Event), idx: uuid.New()}
}

func (s *Storage) UpdateEvent(id uuid.UUID, e storage.Event, ctx context.Context) (storage.Event, error) {

	err := e.Validate()
	if err != nil {
		return storage.Event{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[id]
	if !ok {
		return storage.Event{}, appErrors.ErrIdDoesNotExist
	}
	e.ID = id
	s.data[id] = e

	return e, nil
}

func (s *Storage) ReadEvents(context.Context) ([]storage.Event, error) {

	result := make([]storage.Event, 0)
	s.mu.RLock()
	for _, event := range s.data {
		result = append(result, event)
	}
	s.mu.RUnlock()
	return result, nil
}

func (s *Storage) DeleteEvent(id uuid.UUID, ctx context.Context) (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[id]
	if !ok {
		return uuid.Nil, appErrors.ErrIdDoesNotExist
	}
	delete(s.data, id)
	return id, nil
}

func (s *Storage) CreateEvent(e storage.Event, ctx context.Context) (storage.Event, error) {
	err := e.Validate()
	if err != nil {
		return storage.Event{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e.ID = s.idx
	s.data[s.idx] = e
	s.idx = uuid.New()

	return e, nil

}

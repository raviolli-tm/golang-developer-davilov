package app

import (
	"context"
	"fmt"
	internalhttp "github.com/davilov/hw12_13_14_15_calendar/internal/server/http"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	internalhttp.Logger
}

type Storage interface {
	CreateEvent(e storage.Event) error
	ReadEvents() ([]storage.Event, error)
	UpdateEvent(id uuid.UUID, e storage.Event) error
	DeleteEvent(id uuid.UUID) error
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id uuid.UUID, title string) error {

	err := a.storage.CreateEvent(storage.Event{ID: id, Title: title})
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

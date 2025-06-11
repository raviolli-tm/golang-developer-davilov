package app

import (
	"context"
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
	Error(msg string)
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

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {

	err := a.storage.CreateEvent(event)
	if err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

func (a *App) ReadEvents(ctx context.Context) ([]storage.Event, error) {

	events, err := a.storage.ReadEvents()
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return events, nil
}
func (a *App) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) error {

	err := a.storage.UpdateEvent(id, event)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	return nil
}
func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {

	err := a.storage.DeleteEvent(id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	return nil
}

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
	CreateEvent(e storage.Event, ctx context.Context) (storage.Event, error)
	ReadEvents(ctx context.Context) ([]storage.Event, error)
	UpdateEvent(id uuid.UUID, e storage.Event, ctx context.Context) (storage.Event, error)
	DeleteEvent(id uuid.UUID, ctx context.Context) (uuid.UUID, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {

	event, err := a.storage.CreateEvent(event, ctx)
	if err != nil {
		return storage.Event{}, fmt.Errorf("create event: %w", err)
	}
	return event, nil
}

func (a *App) ReadEvents(ctx context.Context) ([]storage.Event, error) {

	events, err := a.storage.ReadEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return events, nil
}
func (a *App) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) (storage.Event, error) {

	event, err := a.storage.UpdateEvent(id, event, ctx)
	if err != nil {
		return storage.Event{}, fmt.Errorf("update event: %w", err)
	}
	return event, nil
}
func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {

	_, err := a.storage.DeleteEvent(id, ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("delete event: %w", err)
	}
	return id, nil
}

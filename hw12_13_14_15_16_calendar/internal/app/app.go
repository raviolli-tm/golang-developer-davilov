package app

import (
	"context"
	"fmt"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/metrics"
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
	ReadEventsByPeriod(from, to time.Time, ctx context.Context) ([]storage.Event, error)
	UpdateEvent(id uuid.UUID, e storage.Event, ctx context.Context) (storage.Event, error)
	DeleteEvent(id uuid.UUID, ctx context.Context) (uuid.UUID, error)
	Close(ctx context.Context) error
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	event, err := a.storage.CreateEvent(event, ctx)
	if err != nil {
		metrics.EventErrorsTotal.WithLabelValues("create").Inc()
		return storage.Event{}, fmt.Errorf("create event: %w", err)
	}

	metrics.EventsCreatedTotal.Inc()

	return event, nil
}

func (a *App) ReadEvents(ctx context.Context) ([]storage.Event, error) {
	events, err := a.storage.ReadEvents(ctx)
	if err != nil {
		metrics.EventErrorsTotal.WithLabelValues("list").Inc()
		return nil, fmt.Errorf("get events: %w", err)
	}

	return events, nil
}

func (a *App) ReadEventsByPeriod(ctx context.Context, from, to time.Time) ([]storage.Event, error) {
	events, err := a.storage.ReadEventsByPeriod(from, to, ctx)
	if err != nil {
		metrics.EventErrorsTotal.WithLabelValues("list_period").Inc()
		return nil, fmt.Errorf("get events by period: %w", err)
	}

	return events, nil
}

func (a *App) UpdateEvent(ctx context.Context, id uuid.UUID, event storage.Event) (storage.Event, error) {
	event, err := a.storage.UpdateEvent(id, event, ctx)
	if err != nil {
		metrics.EventErrorsTotal.WithLabelValues("update").Inc()
		return storage.Event{}, fmt.Errorf("update event: %w", err)
	}

	metrics.EventsUpdatedTotal.Inc()

	return event, nil
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	_, err := a.storage.DeleteEvent(id, ctx)
	if err != nil {
		metrics.EventErrorsTotal.WithLabelValues("delete").Inc()
		return uuid.Nil, fmt.Errorf("delete event: %w", err)
	}

	metrics.EventsDeletedTotal.Inc()

	return id, nil
}

func (a *App) Close(ctx context.Context) error {
	return a.storage.Close(ctx)
}

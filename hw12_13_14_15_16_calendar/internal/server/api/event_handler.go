package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	genApi "github.com/davilov/hw12_13_14_15_calendar/api/go"
	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Application interface {
	CreateEvent(context.Context, storage.Event) (storage.Event, error)
	ReadEvents(ctx context.Context) ([]storage.Event, error)
	ReadEventsByPeriod(ctx context.Context, from, to time.Time) ([]storage.Event, error)
	DeleteEvent(context.Context, uuid.UUID) (uuid.UUID, error)
	UpdateEvent(context.Context, uuid.UUID, storage.Event) (storage.Event, error)
}

type EventAPIService struct {
	App Application
}

func NewEventAPIService(app Application) *EventAPIService {
	return &EventAPIService{App: app}
}

func (s *EventAPIService) GetCalendarEvents(ctx context.Context) (genApi.ImplResponse, error) {
	events, err := s.App.ReadEvents(ctx)
	if err != nil {
		return errorResponse(err), nil
	}

	return genApi.Response(http.StatusOK, eventsToAPI(events)), nil
}

// GetCalendarEventsByPeriod returns events for the day/week/month that
// contains date.
func (s *EventAPIService) GetCalendarEventsByPeriod(
	ctx context.Context,
	date time.Time,
	period storage.Period,
) (genApi.ImplResponse, error) {
	from, to, err := storage.PeriodRange(date, period)
	if err != nil {
		return genApi.Response(http.StatusUnprocessableEntity, genApi.Error{ //nolint:nilerr
			Code:    "invalid_period",
			Message: err.Error(),
		}), nil
	}

	events, err := s.App.ReadEventsByPeriod(ctx, from, to)
	if err != nil {
		return errorResponse(err), nil
	}

	return genApi.Response(http.StatusOK, eventsToAPI(events)), nil
}

func (s *EventAPIService) AddCalendarEvent(ctx context.Context, event genApi.Event) (genApi.ImplResponse, error) {
	createdEvent, err := s.App.CreateEvent(ctx, apiEventToStorage(event))
	if err != nil {
		return errorResponse(err), nil
	}

	return genApi.Response(http.StatusOK, storageEventToAPI(createdEvent)), nil
}

func (s *EventAPIService) UpdateCalendarEventById(
	ctx context.Context,
	id string,
	event genApi.Event,
) (genApi.ImplResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return genApi.Response(http.StatusUnprocessableEntity, genApi.Error{ //nolint:nilerr
			Code:    "invalid_id",
			Message: "event id must be a valid UUID",
		}), nil
	}

	updatedEvent, err := s.App.UpdateEvent(ctx, uid, apiEventToStorage(event))
	if err != nil {
		return errorResponse(err), nil
	}

	return genApi.Response(http.StatusOK, storageEventToAPI(updatedEvent)), nil
}

func (s *EventAPIService) DeleteCalendarEventById(ctx context.Context, id string) (genApi.ImplResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return genApi.Response(http.StatusUnprocessableEntity, genApi.Error{ //nolint:nilerr
			Code:    "invalid_id",
			Message: "event id must be a valid UUID",
		}), nil
	}

	deletedUUID, err := s.App.DeleteEvent(ctx, uid)
	if err != nil {
		return errorResponse(err), nil
	}

	return genApi.Response(http.StatusOK, map[string]string{"event_id": deletedUUID.String()}), nil
}

// errorResponse maps application errors to HTTP responses. Business errors
// must not be returned as the second value, otherwise the generated error
// handler would replace the response body.
func errorResponse(err error) genApi.ImplResponse {
	status := http.StatusInternalServerError
	code := "internal_error"

	switch {
	case errors.Is(err, appErrors.ErrValidation):
		status = http.StatusUnprocessableEntity
		code = "validation_error"
	case errors.Is(err, appErrors.ErrIdDoesNotExist):
		status = http.StatusNotFound
		code = "not_found"
	}

	return genApi.Response(status, genApi.Error{Code: code, Message: err.Error()})
}

func storageEventToAPI(event storage.Event) genApi.Event {
	return genApi.Event{
		EventId:         event.ID.String(),
		Title:           event.Title,
		DateStart:       event.DateStart,
		DateEnd:         event.DateEnd,
		Description:     event.Description,
		UserId:          int32(event.UserId),
		EventNotifyTime: int32(event.EventNotifyTime),
	}
}

func eventsToAPI(events []storage.Event) []genApi.Event {
	result := make([]genApi.Event, 0, len(events))
	for _, event := range events {
		result = append(result, storageEventToAPI(event))
	}

	return result
}

func apiEventToStorage(event genApi.Event) storage.Event {
	result := storage.Event{
		Title:           event.Title,
		DateStart:       event.DateStart,
		DateEnd:         event.DateEnd,
		Description:     event.Description,
		UserId:          int(event.UserId),
		EventNotifyTime: int(event.EventNotifyTime),
	}

	if event.EventId != "" {
		if id, err := uuid.Parse(event.EventId); err == nil {
			result.ID = id
		}
	}

	return result
}

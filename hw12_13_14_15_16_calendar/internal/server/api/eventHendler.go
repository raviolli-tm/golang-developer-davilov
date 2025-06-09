package api

import (
	"context"
	api "github.com/davilov/hw12_13_14_15_calendar/api/go"
	internalhttp "github.com/davilov/hw12_13_14_15_calendar/internal/server/http"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type EventAPIService struct {
	App internalhttp.Application
}

func NewEventAPIService(app internalhttp.Application) *EventAPIService {
	return &EventAPIService{App: app}
}

func (s *EventAPIService) GetCalendarEvents(ctx context.Context) (api.ImplResponse, error) {
	events, err := s.App.ReadEvents(ctx)
	if err != nil {
		return api.Response(500, err), err
	}

	return api.Response(200, events), nil

}
func (s *EventAPIService) AddCalendarEvent(ctx context.Context, event api.Event) (api.ImplResponse, error) {
	err := s.App.CreateEvent(ctx, s.ApiEventToStorageEvent(event))
	if err != nil {
		return api.Response(500, err), err
	}
	return api.Response(200, "success"), nil
}
func (s *EventAPIService) UpdateCalendarEventById(ctx context.Context, id string, event api.Event) (api.ImplResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return api.Response(422, err), err
	}
	err = s.App.UpdateEvent(ctx, uid, s.ApiEventToStorageEvent(event))
	if err != nil {
		return api.Response(500, err), err
	}
	return api.Response(200, "success"), nil
}
func (s *EventAPIService) DeleteCalendarEventById(ctx context.Context, id string) (api.ImplResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return api.Response(422, err), err
	}
	err = s.App.DeleteEvent(ctx, uid)
	if err != nil {
		return api.Response(500, err), err
	}
	return api.Response(200, "success"), nil
}

func (s *EventAPIService) StorageEventToApiEvent(event storage.Event) api.Event {
	resultEvent := api.Event{}
	resultEvent.ID = event.ID.String()
	resultEvent.Title = event.Title
	resultEvent.UserId = int32(event.UserId)
	resultEvent.DateStart = event.DateStart
	resultEvent.DateEnd = event.DateEnd
	resultEvent.Description = event.Description
	resultEvent.EventNotifyTime = int32(event.EventNotifyTime)

	return resultEvent

}

func (s *EventAPIService) ApiEventToStorageEvent(event api.Event) storage.Event {

	resultEvent := storage.Event{}
	resultEvent.ID, _ = uuid.Parse(event.ID)
	resultEvent.Title = event.Title
	resultEvent.UserId = int(event.UserId)
	resultEvent.DateStart = event.DateStart
	resultEvent.DateEnd = event.DateEnd
	resultEvent.Description = event.Description
	resultEvent.EventNotifyTime = int(event.EventNotifyTime)

	return resultEvent
}

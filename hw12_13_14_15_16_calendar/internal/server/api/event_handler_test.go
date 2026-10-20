package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	genApi "github.com/davilov/hw12_13_14_15_calendar/api/go"
	"github.com/davilov/hw12_13_14_15_calendar/appErrors"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeApp struct {
	events   []storage.Event
	listFrom time.Time
	listTo   time.Time
}

func (f *fakeApp) CreateEvent(_ context.Context, event storage.Event) (storage.Event, error) {
	if err := event.Validate(); err != nil {
		return storage.Event{}, err
	}

	event.ID = uuid.New()
	f.events = append(f.events, event)

	return event, nil
}

func (f *fakeApp) ReadEvents(_ context.Context) ([]storage.Event, error) {
	return f.events, nil
}

func (f *fakeApp) ReadEventsByPeriod(_ context.Context, from, to time.Time) ([]storage.Event, error) {
	f.listFrom, f.listTo = from, to

	result := make([]storage.Event, 0)
	for _, event := range f.events {
		if !event.DateStart.Before(from) && event.DateStart.Before(to) {
			result = append(result, event)
		}
	}

	return result, nil
}

func (f *fakeApp) UpdateEvent(_ context.Context, id uuid.UUID, event storage.Event) (storage.Event, error) {
	for i, item := range f.events {
		if item.ID == id {
			event.ID = id
			f.events[i] = event
			return event, nil
		}
	}

	return storage.Event{}, appErrors.ErrIdDoesNotExist
}

func (f *fakeApp) DeleteEvent(_ context.Context, id uuid.UUID) (uuid.UUID, error) {
	for i, item := range f.events {
		if item.ID == id {
			f.events = append(f.events[:i], f.events[i+1:]...)
			return id, nil
		}
	}

	return uuid.Nil, appErrors.ErrIdDoesNotExist
}

func newTestServer(app Application) *httptest.Server {
	mux := http.NewServeMux()
	service := NewEventAPIService(app)
	controller := genApi.NewCalendarEventsAPIController(service)
	router := genApi.NewRouter(NewPeriodRouter(service), controller)
	mux.Handle("/calendar", router)
	mux.Handle("/calendar/", router)

	return httptest.NewServer(mux)
}

func doRequest(
	t *testing.T,
	server *httptest.Server,
	method, path string,
	body []byte,
) *http.Response {
	t.Helper()

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, server.URL+path, reader)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := server.Client().Do(req)
	require.NoError(t, err)

	return resp
}

func TestEventHandler(t *testing.T) {
	app := &fakeApp{}
	server := newTestServer(app)
	defer server.Close()

	start := time.Date(2026, 10, 20, 10, 0, 0, 0, time.UTC)
	payload := genApi.Event{
		Title:           "test event",
		DateStart:       start,
		DateEnd:         start.Add(time.Hour),
		Description:     "description",
		UserId:          1,
		EventNotifyTime: 15,
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	resp := doRequest(t, server, http.MethodPost, "/calendar", body)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var created genApi.Event
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	require.NotEmpty(t, created.EventId)

	resp = doRequest(t, server, http.MethodGet, "/calendar", nil)
	defer resp.Body.Close()

	var events []genApi.Event
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&events))
	require.Len(t, events, 1)
	require.Equal(t, created.EventId, events[0].EventId)
}

func TestEventHandlerValidation(t *testing.T) {
	app := &fakeApp{}
	server := newTestServer(app)
	defer server.Close()

	body, err := json.Marshal(genApi.Event{Title: ""})
	require.NoError(t, err)

	resp := doRequest(t, server, http.MethodPost, "/calendar", body)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var apiError genApi.Error
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiError))
	require.Equal(t, "validation_error", apiError.Code)
}

func TestEventHandlerPeriod(t *testing.T) {
	start := time.Date(2026, 10, 20, 10, 0, 0, 0, time.UTC)
	app := &fakeApp{events: []storage.Event{{
		ID:        uuid.New(),
		Title:     "period event",
		DateStart: start,
		DateEnd:   start.Add(time.Hour),
		UserId:    1,
	}}}
	server := newTestServer(app)
	defer server.Close()

	resp := doRequest(t, server, http.MethodGet, "/calendar/day?date=2026-10-20", nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC), app.listFrom)
	require.Equal(t, time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC), app.listTo)

	var events []genApi.Event
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&events))
	require.Len(t, events, 1)
}

func TestEventHandlerInvalidDate(t *testing.T) {
	server := newTestServer(&fakeApp{})
	defer server.Close()

	resp := doRequest(t, server, http.MethodGet, "/calendar/week?date=invalid", nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestEventHandlerNotFound(t *testing.T) {
	server := newTestServer(&fakeApp{})
	defer server.Close()

	resp := doRequest(t, server, http.MethodDelete, "/calendar/"+uuid.New().String(), nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

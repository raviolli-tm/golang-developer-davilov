package internalhttp

import (
	"context"
	genApi "github.com/davilov/hw12_13_14_15_calendar/api/go"
	"github.com/davilov/hw12_13_14_15_calendar/internal/server/api"
	"github.com/davilov/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Server struct {
	Logger      Logger
	Application Application
	ctx         context.Context
}

type Logger interface {
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(context.Context, storage.Event) error
	ReadEvents(ctx context.Context) ([]storage.Event, error)
	DeleteEvent(context.Context, uuid.UUID) error
	UpdateEvent(context.Context, uuid.UUID, storage.Event) error
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{Logger: logger, Application: app}
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx = ctx

	handler := http.NewServeMux()
	handler.HandleFunc("/", s.homeHandler)
	handler.HandleFunc("/hello-world", s.helloWorld)

	calendarApi := api.NewEventAPIService(s.Application)
	controllerEvent := genApi.NewCalendarEventsAPIController(calendarApi)

	muxRouter := genApi.NewRouter(controllerEvent)

	handler.Handle("/events", muxRouter)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(handler, s.Logger),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()

	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

}

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {

	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// TODO

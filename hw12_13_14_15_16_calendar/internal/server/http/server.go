package internalhttp

import (
	"context"
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
	CreateEvent(ctx context.Context, id uuid.UUID, title string) error
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{Logger: logger, Application: app}
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx = ctx
	handler := http.NewServeMux()
	handler.Handle("/", loggingMiddleware(s.homeHandler, s.Logger))
	handler.Handle("/hello-world", loggingMiddleware(s.helloWorld, s.Logger))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
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

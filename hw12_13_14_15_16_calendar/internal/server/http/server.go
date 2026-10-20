package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	Logger      app.Logger
	Application *app.App
	addr        string
	httpServer  *http.Server
}

func NewServer(logger app.Logger, application *app.App, addr string) *Server {
	return &Server{Logger: logger, Application: application, addr: addr}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(http.HandlerFunc(s.homeHandler), s.Logger))
	mux.Handle("/hello-world", loggingMiddleware(http.HandlerFunc(s.helloWorld), s.Logger))

	s.httpServer = &http.Server{
		Addr:              s.addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Stop(shutdownCtx); err != nil {
			s.Logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	s.Logger.Info("calendar http server is listening on " + s.addr)

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("Hello World")); err != nil {
		s.Logger.Error("failed to write hello response: " + err.Error())
	}
}

package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	genApi "github.com/davilov/hw12_13_14_15_calendar/api/go"
	"github.com/davilov/hw12_13_14_15_calendar/internal/app"
	"github.com/davilov/hw12_13_14_15_calendar/internal/metrics"
	"github.com/davilov/hw12_13_14_15_calendar/internal/server/api"
)

type Server struct {
	logger      app.Logger
	application api.Application
	addr        string
	httpServer  *http.Server
}

func NewServer(logger app.Logger, application api.Application, addr string) *Server {
	return &Server{logger: logger, application: application, addr: addr}
}

func (s *Server) Start(ctx context.Context) error {
	handler := http.NewServeMux()
	handler.HandleFunc("/", s.helloWorld)
	handler.HandleFunc("/hello", s.helloWorld)
	handler.HandleFunc("/hello-world", s.helloWorld)
	handler.Handle("/metrics", metrics.Handler())

	calendarAPI := api.NewEventAPIService(s.application)
	controller := genApi.NewCalendarEventsAPIController(calendarAPI)
	router := genApi.NewRouter(api.NewPeriodRouter(calendarAPI), controller)
	handler.Handle("/calendar", router)
	handler.Handle("/calendar/", router)

	s.httpServer = &http.Server{
		Addr:              s.addr,
		Handler:           loggingMiddleware(metricsMiddleware(handler), s.logger),
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
			s.logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	s.logger.Info(fmt.Sprintf("calendar http server is listening on %s", s.addr))

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

func (s *Server) helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("Hello World")); err != nil {
		s.logger.Error("failed to write hello response: " + err.Error())
	}
}

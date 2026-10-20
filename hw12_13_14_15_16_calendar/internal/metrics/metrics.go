// Package metrics contains the Prometheus metrics of the calendar service.
//
// Metrics are registered in the default Prometheus registry, so the /metrics
// endpoint of the API and the dedicated metrics servers of the background
// processes expose the standard Go/process collectors as well.
package metrics

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTPRequestsTotal counts HTTP requests by method, route and status.
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	// HTTPRequestDuration observes HTTP request latency.
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "calendar",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request latency in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path"})

	// HTTPRequestsInFlight tracks the number of requests being processed.
	HTTPRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "calendar",
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "Number of HTTP requests currently being processed.",
	})

	// EventsCreatedTotal counts successfully created events.
	EventsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "events",
		Name:      "created_total",
		Help:      "Total number of created calendar events.",
	})

	// EventsUpdatedTotal counts successfully updated events.
	EventsUpdatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "events",
		Name:      "updated_total",
		Help:      "Total number of updated calendar events.",
	})

	// EventsDeletedTotal counts successfully deleted events.
	EventsDeletedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "events",
		Name:      "deleted_total",
		Help:      "Total number of deleted calendar events.",
	})

	// EventErrorsTotal counts failed event operations by operation name.
	EventErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "events",
		Name:      "errors_total",
		Help:      "Total number of failed event operations.",
	}, []string{"operation"})

	// NotificationsStoredTotal counts notifications persisted by the storer.
	NotificationsStoredTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "storer",
		Name:      "notifications_stored_total",
		Help:      "Total number of notifications saved to the database.",
	})

	// NotificationsErrorsTotal counts notification processing errors.
	NotificationsErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "storer",
		Name:      "notifications_errors_total",
		Help:      "Total number of notification processing errors.",
	})

	// SchedulerScansTotal counts scheduler scan iterations.
	SchedulerScansTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "scheduler",
		Name:      "scans_total",
		Help:      "Total number of scheduler scan iterations.",
	})

	// SchedulerErrorsTotal counts scheduler errors.
	SchedulerErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "scheduler",
		Name:      "errors_total",
		Help:      "Total number of scheduler errors.",
	})

	// SchedulerNotificationsSentTotal counts notifications published to Kafka.
	SchedulerNotificationsSentTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "scheduler",
		Name:      "notifications_sent_total",
		Help:      "Total number of notifications published to the broker.",
	})

	// SchedulerEventsDeletedTotal counts stale events removed by the scheduler.
	SchedulerEventsDeletedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calendar",
		Subsystem: "scheduler",
		Name:      "events_deleted_total",
		Help:      "Total number of stale events deleted by the scheduler.",
	})
)

// Handler returns the Prometheus HTTP handler.
func Handler() http.Handler {
	return promhttp.Handler()
}

// Serve starts an HTTP server exposing /metrics on addr and shuts it down
// when ctx is cancelled.
func Serve(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", Handler())

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

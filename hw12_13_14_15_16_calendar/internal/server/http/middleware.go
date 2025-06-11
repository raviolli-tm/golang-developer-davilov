package internalhttp

import (
	"fmt"
	"github.com/davilov/hw12_13_14_15_calendar/internal/app"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler, logger app.Logger) http.HandlerFunc { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		startTime := time.Now()
		next.ServeHTTP(lrw, r)
		duration := startTime.Sub(time.Now()).Milliseconds()
		statusCode := lrw.statusCode

		if statusCode >= 200 && statusCode <= 299 {
			logger.Info(
				fmt.Sprintf("%s [%s] %s %s %s %d %dms",
					r.RemoteAddr,
					time.Now().Format(time.RFC3339),
					r.Method,
					r.URL.Path,
					r.Proto,
					statusCode,
					duration))

		} else {
			logger.Info(
				fmt.Sprintf("%s [%s] %s %s %s %d %dms",
					r.RemoteAddr,
					time.Now().Format(time.RFC3339),
					r.Method,
					r.URL.Path,
					r.Proto,
					statusCode,
					duration))
		}
	}
}

package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/app"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	return lrw.ResponseWriter.Write(data)
}

func loggingMiddleware(next http.Handler, logger app.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		startTime := time.Now()
		next.ServeHTTP(lrw, r)
		duration := time.Since(startTime)

		clientIP := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			clientIP = host
		}

		message := fmt.Sprintf(
			`%s [%s] "%s %s %s" %d %s %q`,
			clientIP,
			startTime.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.RequestURI(),
			r.Proto,
			lrw.statusCode,
			duration.Round(time.Microsecond),
			r.UserAgent(),
		)

		if lrw.statusCode >= http.StatusInternalServerError {
			logger.Error(message)
			return
		}

		logger.Info(message)
	}
}

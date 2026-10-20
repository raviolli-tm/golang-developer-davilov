package internalhttp

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/davilov/hw12_13_14_15_calendar/internal/app"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}

	n, err := lrw.ResponseWriter.Write(data)
	lrw.written += int64(n)

	return n, err
}

// Flush keeps streaming endpoints working through the middleware.
func (lrw *loggingResponseWriter) Flush() {
	if flusher, ok := lrw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack keeps websocket-style upgrades working through the middleware.
func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("response writer does not support hijacking")
	}

	return hijacker.Hijack()
}

func loggingMiddleware(next http.Handler, logger app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		} else {
			logger.Info(message)
		}
	})
}

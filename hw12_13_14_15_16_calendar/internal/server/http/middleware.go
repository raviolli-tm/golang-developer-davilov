package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func loggingMiddleware(next http.HandlerFunc, logger Logger) http.HandlerFunc { //nolint:unused
	return func(w http.ResponseWriter, r *http.Request) {

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		startTime := time.Now()
		next(lrw, r)
		duration := startTime.Sub(time.Now()).Milliseconds()
		lrw.Header().Get("Content-Type")
		lrw.Header().Get("X-Forwarded-For")
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

		//IP клиента;
		//дата и время запроса;
		//метод, path и версия HTTP;
		//код ответа;
		//latency (время обработки запроса, посчитанное, например, с помощью middleware);
		//user agent, если есть.

		//66.249.65.3 [25/Feb/2020:19:11:24 +0600] GET /hello?q=1 HTTP/1.1 200 30 "Mozilla/5.0"
	}
}

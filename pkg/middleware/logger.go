package middleware

import (
	strftime "github.com/itchyny/timefmt-go"
	"log"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		logger.Printf("%s,%s,%s,%q,%d,%d",
			strftime.Format(time.Now(), "%Y/%m/%d-%H:%M:%S"),
			r.Method,
			r.RemoteAddr,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start).Microseconds(),
		)
	})
}

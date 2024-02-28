package main

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(logger *log.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		duration := end.Sub(start)
		logger.Printf("%s,%s,%s,%q,%d,%d", time.Now(), r.Method, r.RemoteAddr, r.URL.Path, w.Header().Get("X-Status-Code"), duration.Milliseconds())
	}
}

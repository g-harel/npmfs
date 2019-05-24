package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(statusCode int) {
	s.status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{w, http.StatusOK}
		h.ServeHTTP(rec, r)
		log.Printf("%v %v %v - %vms", r.Method, r.RequestURI, rec.status, int64(time.Since(start)/time.Millisecond))
	})
}

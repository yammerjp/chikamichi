package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func WithLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		method := r.Method
		url := r.URL.String()
		userAgent := r.Header.Get("User-Agent")

		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		h.ServeHTTP(recorder, r)
		status := strconv.Itoa(recorder.Status)
		log.Println(strings.Join([]string{remoteAddr, status, method, url, userAgent}, " "))
	})
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

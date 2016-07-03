package main

import (
	"net/http"
	"net/http/httptest"
)

type CORSMiddleware struct {
	// Handler is the handler to be wrapped
	Handler http.Handler
}

func NewCORSMiddleware(h http.Handler) *CORSMiddleware {
	return &CORSMiddleware{
		Handler: h,
	}
}

func (m *CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()

	m.Handler.ServeHTTP(rec, r)

	// recreate the response
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	// add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

	w.WriteHeader(rec.Code)
	w.Write(rec.Body.Bytes())
}

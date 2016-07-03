package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

func Test_HTTPMiddleware_CORS_Factory(t *testing.T) {
	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	h := &tmHTTPHandler{
		hFn: handlerFn,
	}

	m := NewCORSMiddleware(h)

	ar.NotNil(t, m, "empty element returned")
	ar.IsType(t, &CORSMiddleware{}, m)

	a.Equal(t, h, m.Handler, "Handler is not attached")
}

func Test_HTTPMiddleware_CORS(t *testing.T) {
	h := &tmHTTPHandler{
		hFn: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-A", "123")
			w.Header().Set("X-Test-B", "456")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("0123456789"))
		},
	}

	m := NewCORSMiddleware(h)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	res := httptest.NewRecorder()
	m.ServeHTTP(res, req)

	hExp := map[string]string{
		"X-Test-A":                     "123",
		"X-Test-B":                     "456",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, DELETE, PUT, PATCH, OPTIONS",
		"Access-Control-Allow-Headers": "Origin, Content-Type",
	}
	for hName, hVal := range hExp {
		a.Equal(t, hVal, res.Header().Get(hName), "mismatch on response header: %s", hName)
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/uber-go/zap"
	"github.com/uber-go/zap/spy"

	a "github.com/stretchr/testify/assert"
	ar "github.com/stretchr/testify/require"
)

type tmHTTPHandler struct {
	hFn func(w http.ResponseWriter, r *http.Request)
}

func (h *tmHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.hFn(w, r)
}

func Test_HTTPMiddleware_Factory(t *testing.T) {
	l, _ := spy.New()
	l.SetLevel(zap.DebugLevel)

	handlerFn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	h := &tmHTTPHandler{
		hFn: handlerFn,
	}

	m := NewLoggingMiddleware(h, l)

	ar.NotNil(t, m, "empty element returned")
	ar.IsType(t, &LoggingMiddleware{}, m)

	a.NotNil(t, m.TimeNow, "TimeNow not initialised")
	a.Equal(t, h, m.Handler, "Handler is not attached")
	a.Equal(t, l, m.Logger, "Logger is not attached")
}

func Test_HTTPMiddleware_Logging(t *testing.T) {
	// fake duration of each request processing
	timeFakeCh := make(chan time.Time, 2)
	timePreRequest := time.Date(2016, time.May, 29, 10, 11, 12, 13, time.UTC)
	timePostRequest := timePreRequest.Add(time.Millisecond * 3)
	timeDurationMS := timePostRequest.Sub(timePreRequest).Seconds() * 1e-3

	tests := map[string]struct {
		method    string
		reqBody   string
		handlerFn func(w http.ResponseWriter, r *http.Request)
		headers   map[string]string
		exp       spy.Log
	}{
		"success, GET, ok": {
			method:  http.MethodGet,
			headers: map[string]string{"X-Test-A": "123", "X-Test-B": "456"},
			handlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Test-A", "123")
				w.Header().Set("X-Test-B", "456")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("0123456789"))
			},
			exp: spy.Log{
				Level: zap.InfoLevel,
				Msg:   "request:done",
				Fields: []zap.Field{
					zap.String("req:method", http.MethodGet),
					zap.String("req:proto", "HTTP/1.1"),
					zap.String("req:host", "example.com"),
					zap.String("req:URI", "/foo"),
					zap.Int64("req:contentLength", 0),
					zap.Int("res:status", http.StatusOK),
					zap.Int("res:contentLength", 10),
					zap.Float64("req:duration:ms", timeDurationMS),
				},
			},
		},
		"failure, POST, BadReq, empty return": {
			method:  http.MethodPost,
			reqBody: "01234567890123456789",
			handlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				return
			},
			exp: spy.Log{
				Level: zap.InfoLevel,
				Msg:   "request:done",
				Fields: []zap.Field{
					zap.String("req:method", http.MethodPost),
					zap.String("req:proto", "HTTP/1.1"),
					zap.String("req:host", "example.com"),
					zap.String("req:URI", "/foo"),
					zap.Int64("req:contentLength", 20),
					zap.Int("res:status", http.StatusBadRequest),
					zap.Int("res:contentLength", 0),
					zap.Float64("req:duration:ms", timeDurationMS),
				},
			},
		},
		"failure, GET, ServiceUnavailable": {
			method: http.MethodGet,
			handlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("0123456789"))
				return
			},
			exp: spy.Log{
				Level: zap.WarnLevel,
				Msg:   "request:done",
				Fields: []zap.Field{
					zap.String("req:method", http.MethodGet),
					zap.String("req:proto", "HTTP/1.1"),
					zap.String("req:host", "example.com"),
					zap.String("req:URI", "/foo"),
					zap.Int64("req:contentLength", 0),
					zap.Int("res:status", http.StatusServiceUnavailable),
					zap.Int("res:contentLength", 10),
					zap.Float64("req:duration:ms", timeDurationMS),
				},
			},
		},
		"failure, GET, InternalServerError": {
			method: http.MethodGet,
			handlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				return
			},
			exp: spy.Log{
				Level: zap.ErrorLevel,
				Msg:   "request:done",
				Fields: []zap.Field{
					zap.String("req:method", http.MethodGet),
					zap.String("req:proto", "HTTP/1.1"),
					zap.String("req:host", "example.com"),
					zap.String("req:URI", "/foo"),
					zap.Int64("req:contentLength", 0),
					zap.Int("res:status", http.StatusInternalServerError),
					zap.Int("res:contentLength", 0),
					zap.Float64("req:duration:ms", timeDurationMS),
				},
			},
		},
	}

	for sym, tc := range tests {
		lgr, sink := spy.New()
		lgr.SetLevel(zap.DebugLevel)

		timeFakeCh <- timePreRequest
		timeFakeCh <- timePostRequest

		h := &tmHTTPHandler{
			hFn: tc.handlerFn,
		}

		m := LoggingMiddleware{
			Logger:  lgr,
			Handler: h,
			TimeNow: func() time.Time {
				return <-timeFakeCh
			},
		}
		reqBodyReader := strings.NewReader(tc.reqBody)
		req, _ := http.NewRequest(tc.method, "http://example.com/foo", reqBodyReader)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)

		for hName, hVal := range tc.headers {
			a.Equal(t, hVal, res.Header().Get(hName), "[%s] mismatch on response header: %s", sym, hName)
		}

		got := sink.Logs()
		if a.Len(t, got, 1, "[%s] incorrect number of logs generated", sym) {
			a.Equal(t, tc.exp, got[0], "[%s] Unexpected output from sampled logger.", sym)
		}
	}
}

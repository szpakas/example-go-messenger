package main

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/uber-go/zap"
)

// LoggingMiddleware provides HTTP middleware which allows logging of the request and response.
type LoggingMiddleware struct {
	// Handler is the handler to be wrapped
	Handler http.Handler

	// Logger is the instance of zap.Logger used in logging
	Logger zap.Logger

	// TimeNow is testing helper for time sensitive tests. It defaults to time.Now function.
	TimeNow func() time.Time
}

func NewLoggingMiddleware(h http.Handler, l zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		Handler: h,
		Logger:  l,
		TimeNow: time.Now,
	}
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqStartedTime := m.TimeNow()

	rec := httptest.NewRecorder()

	m.Handler.ServeHTTP(rec, r)

	// -- recreate the response
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	w.WriteHeader(rec.Code)
	w.Write(rec.Body.Bytes())

	// -- log
	var ll zap.Level
	switch rec.Code {
	case http.StatusInternalServerError:
		ll = zap.ErrorLevel
	case http.StatusServiceUnavailable:
		ll = zap.WarnLevel
	default:
		ll = zap.InfoLevel
	}
	m.Logger.Log(
		ll,
		"request:done",
		zap.String("req:method", r.Method),
		zap.String("req:proto", r.Proto),
		zap.String("req:host", r.Host),
		zap.String("req:URI", r.URL.Path),
		zap.Int64("req:contentLength", r.ContentLength),
		zap.Int("res:status", rec.Code),
		zap.Int("res:contentLength", rec.Body.Len()),
		zap.Float64("req:duration:ms", m.TimeNow().Sub(reqStartedTime).Seconds()*1e-3),
	)
}

package main

//noinspection SpellCheckingInspection
import (
	"fmt"

	"github.com/uber-go/zap"
	"github.com/vrischmann/envconfig"
)

const (
	// ConfigAppPrefix prefixes all ENV values used to config the program.
	ConfigAppPrefix = "APP"
)

type config struct {
	// HTTPHost is address on which HTTP server endpoint is listening.
	HTTPHost string `envconfig:"default=0.0.0.0"`

	// HTTPPort is a port number on which HTTP server endpoint is listening.
	HTTPPort int `envconfig:"default=8080"`

	// LogLevel is a minimal log severity required for the message to be logged.
	// Valid levels: [debug, info, warn, error, fatal, panic, none].
	LogLevel string `envconfig:"default=info"`
}

func main() {
	lgr := zap.NewJSON()

	// - config from env
	cfg := &config{}
	if err := envconfig.InitWithPrefix(&cfg, ConfigAppPrefix); err != nil {
		lgr.Fatal(err.Error())
	}

	// -- logging
	var logLevel zap.Level
	if err := logLevel.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		lgr.Fatal(err.Error())
	}

	lgr.SetLevel(logLevel)
	lgr.Debug(fmt.Sprintf("Parsed config from env => %+v", *cfg))

	lgr.Info("starting")

	st := NewMemoryStorage()
	h := NewHTTPDefaultHandler(st)
	ml := NewLoggingMiddleware(h, lgr)
	s := NewHTTPServer(cfg.HTTPHost, cfg.HTTPPort, ml)

	if err := s.ListenAndServe(); err != nil {
		lgr.Fatal(err.Error())
	}
}

package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"log/slog"

	"github.com/mcosta74/wingetsvc"
	"github.com/oklog/run"
)

func main() {
	config := wingetsvc.LoadConfig(os.Args[1:])
	logger := setupLogger(config)

	logger.Info("Service Started")
	defer logger.Info("Service Stopped")

	var (
		controller  = wingetsvc.NewWingetController(logger)
		service     = wingetsvc.NewService(controller)
		endpoints   = wingetsvc.MakeEndpoints(service)
		httpHandler = wingetsvc.MakeHTTPHandler(endpoints, logger.With("transport", "HTTP"))
	)

	var g run.Group
	{
		// Signal Handler
		g.Add(run.SignalHandler(context.Background(), syscall.SIGTERM, syscall.SIGINT))
	}

	{
		// HTTP Handler
		httpLogger := logger.With("transport", "HTTP")

		listener, err := net.Listen("tcp", config.HTTPAddr)
		if err != nil {
			httpLogger.Error("Failed to listen", "err", err)
			os.Exit(1)
		}

		g.Add(func() error {
			httpLogger.Info("Listening", "address", listener.Addr())
			return http.Serve(listener, httpHandler)
		}, func(err error) {
			listener.Close()
		})
	}
	logger.Info("group stopped", "err", g.Run())
}

func setupLogger(config *wingetsvc.Config) *slog.Logger {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			// print time in UTC
			a.Value = slog.TimeValue(a.Value.Time().UTC())
		}
		if a.Key == slog.SourceKey {
			// do not print directories
			a.Value = slog.StringValue(filepath.Base(a.Value.String()))
		}
		return a
	}

	opts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       config.LogLevel,
		ReplaceAttr: replace,
	}

	var handler slog.Handler = slog.NewJSONHandler(os.Stderr, opts)
	if config.LogFormat == "text" {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}

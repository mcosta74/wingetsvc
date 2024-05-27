package wingetsvc

import (
	"flag"
	"fmt"
	"os"

	"log/slog"
)

type Config struct {
	LogLevel  slog.Level `json:"log_level,omitempty"`
	LogFormat string     `json:"log_format,omitempty"`
	HTTPAddr  string     `json:"http_addr,omitempty"`
}

func LoadConfig(args []string) *Config {
	config := &Config{
		LogLevel:  slog.LevelInfo,
		LogFormat: "text",
		HTTPAddr:  ":8080",
	}

	loadEnv(config)
	loadFlags(config, args)
	return config
}

func loadEnv(config *Config) {
	if val, ok := loadStringEnv("WGSVC_LOG_LEVEL"); ok {
		config.LogLevel = slogLevelFromString(val, config.LogLevel)
	}

	if val, ok := loadStringEnv("WGSVC_LOG_FORMAT"); ok {
		config.LogFormat = val
	}

	if val, ok := loadStringEnv("WGSVC_HTTP_ADDR"); ok {
		config.HTTPAddr = val
	}
}

func loadStringEnv(name string) (string, bool) {
	return os.LookupEnv(name)
}

func loadFlags(config *Config, args []string) {
	fs := flag.NewFlagSet("wingetsvc", flag.ExitOnError)

	const (
		logLevelName  = "log-level"
		logFormatName = "log-format"
		httpAddrName  = "http-addr"
	)
	fs.String(
		logLevelName, config.LogLevel.String(),
		fmt.Sprintf("Log level (%s, %s, %s, %s)", slog.LevelDebug.String(), slog.LevelInfo.String(), slog.LevelWarn.String(), slog.LevelError.String()),
	)
	fs.String(logFormatName, config.LogFormat, "Log format (text, json)")
	fs.String(httpAddrName, config.HTTPAddr, "HTTP Listen address")
	fs.Parse(args)

	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case logLevelName:
			config.LogLevel = slogLevelFromString(f.Value.String(), config.LogLevel)

		case logFormatName:
			config.LogFormat = f.Value.String()

		case httpAddrName:
			config.HTTPAddr = f.Value.String()
		}
	})
}

func slogLevelFromString(val string, defaultVal slog.Level) slog.Level {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(val)); err == nil {
		return lvl
	}
	return defaultVal
}

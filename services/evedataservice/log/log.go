package log

import (
	"github.com/andrewapj/arcturus/configuration"
	config "github.com/andrewapj/dotenvconfig"
	"log/slog"
	"os"
	"strings"
)

// Configure configures slog to use json logging and the given logging level.
func Configure() {

	level, err := config.GetKey(configuration.LogLevel)
	if err != nil {
		level = ""
	}
	level = strings.ToLower(level)

	var slogLevel slog.Level
	if level == "debug" {
		slogLevel = slog.LevelDebug
	} else if level == "info" {
		slogLevel = slog.LevelInfo
	} else if level == "warn" {
		slogLevel = slog.LevelWarn
	} else if level == "error" {
		slogLevel = slog.LevelError
	} else {
		slogLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))
	logger.Info("created logger with level: " + slogLevel.String())
	slog.SetDefault(logger)
}

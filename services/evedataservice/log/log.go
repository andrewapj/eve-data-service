package log

import (
	"github.com/andrewapj/arcturus/config"
	"log/slog"
	"os"
	"strings"
)

// Configure configures slog to use json logging and the given logging level.
func Configure() {

	level := strings.ToLower(config.LogLevel())

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
	//goland:noinspection GoDfaNilDereference
	logger.Info("created logger", slog.String("level", slogLevel.String()))
	slog.SetDefault(logger)
}

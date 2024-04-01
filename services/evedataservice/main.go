package main

import (
	"embed"
	"github.com/andrewapj/arcturus/config"
	"github.com/andrewapj/arcturus/log"
	"log/slog"
	"os"
)

//go:embed "*.env"
var configFS embed.FS

func main() {

	config.Load(configFS)
	log.Configure()

	slog.Info("application started", "pid", os.Getpid())
}

package main

import (
	"context"
	"embed"
	"github.com/andrewapj/arcturus/config"
	"github.com/andrewapj/arcturus/db"
	"github.com/andrewapj/arcturus/loader"
	"github.com/andrewapj/arcturus/log"
	"log/slog"
	"os"
)

//go:embed "*.env"
var configFS embed.FS

func main() {

	config.Load(configFS)
	log.Configure()

	err := db.Connect(context.Background())
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	loader.Start(context.Background())

	slog.Info("application started", "pid", os.Getpid())
}

package loader

import (
	"context"
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/config"
	"github.com/andrewapj/arcturus/db"
	"github.com/andrewapj/arcturus/esi"
	"log/slog"
	"time"
)

type loader interface {
	load(ctx context.Context, client *esi.Client) (time.Time, error)
	name() string
	lock() string
}

func Start(ctx context.Context) {

	loaders := []loader{
		newStatusLoader(),
	}

	client := esi.NewClient()

	for {
		for _, loader := range loaders {
			load(ctx, &client, loader)
		}

		scheduledTime := clock.GetTime().
			Truncate(time.Duration(config.LoaderIntervalSeconds()) * time.Second).
			Add(time.Duration(config.LoaderIntervalSeconds()) * time.Second)
		waitUntil(scheduledTime)
	}
}

// waitUntil will block until a specific time has been reached.
func waitUntil(targetTime time.Time) {
	now := clock.GetTime()
	if targetTime.After(now) {
		time.Sleep(targetTime.Sub(now))
	}
}

// load is responsible for loading data into the database. It will call a more specialised loader that will
// fetch the data and save it to the database.
func load(ctx context.Context, client *esi.Client, loader loader) {

	slog.Info("starting loader", "loader", loader.name())

	success := db.AcquireLock(ctx, loader.lock())
	if !success {
		slog.Debug("failed to acquire lock", "lock", loader.lock())
		return
	}
	defer func() {
		err := db.ReleaseLock(ctx, loader.lock())
		if err != nil {
			slog.Error("failed to release lock", "lock", loader.lock(), "err", err)
		}
	}()

	expiry, err := db.FindExpiry(ctx, loader.name())
	if err != nil {
		slog.Error("unable to find a valid expiry time during loading", "loader", loader.name(), "err", err)
		return
	}

	if clock.GetTime().After(expiry) {
		slog.Info("updating expired data", "loader", loader.name())
		expires, err := loader.load(ctx, client)
		if err != nil {
			slog.Error("unable to load data", "loader", loader.name(), "err", err)
			return
		}

		err = db.InsertOrUpdateExpiry(ctx, loader.name(), expires)
		if err != nil {
			slog.Error("unable to insert or update expiry", "loader", loader.name(), "err", err)
			return
		}
	}

	slog.Info("completed loader", "loader", loader.name())
}

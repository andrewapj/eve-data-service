package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/andrewapj/arcturus/clock"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

const (
	findExpirySql           = "SELECT expires FROM expiry WHERE kind = $1"
	insertOrUpdateExpirySql = `
	INSERT INTO expiry (kind, expires)
	VALUES ($1, $2) ON CONFLICT (kind) DO UPDATE 
    SET expires = $2;`
)

// FindExpiry will find an expiry for a kind. It returns an error if one does not exist.
func FindExpiry(ctx context.Context, kind string) (time.Time, error) {

	var result time.Time
	err := pool.QueryRow(ctx, findExpirySql, kind).Scan(&result)
	if errors.Is(err, pgx.ErrNoRows) {
		slog.Warn("could not find expiry, returning expired value", "kind", kind)
		return clock.GetTime().Add(-5 * time.Second), nil
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("could not get expiry for kind %s. %w", kind, err)
	}
	return result, nil
}

// InsertOrUpdateExpiry will add a new expiry for a kind, If one already exists it will be overwritten.
func InsertOrUpdateExpiry(ctx context.Context, kind string, expiry time.Time) error {

	tag, err := pool.Exec(ctx, insertOrUpdateExpirySql, kind, expiry)
	if err != nil {
		return fmt.Errorf("could not update expiry for kind %s. %w", kind, err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("could not update expiry for kind %s. RowsAffected %d", kind, tag.RowsAffected())
	}
	return nil
}

package db

import (
	"context"
	"fmt"
)

const (
	acquireLockSql = "INSERT INTO lock (name) VALUES ($1)"
	releaseLockSql = "DELETE FROM lock WHERE name = $1"
)

// AcquireLock attempts to acquire a lock using the given name. If it succeeds it will return true,
// otherwise it returns false.
func AcquireLock(ctx context.Context, name string) bool {

	tag, err := db.Exec(ctx, acquireLockSql, name)
	if err == nil && tag.RowsAffected() == 1 {
		return true

	}
	return false
}

// ReleaseLock attempts to release a lock using the given name. If it is unable to do so or the lock was already
// released it returns an error.
func ReleaseLock(ctx context.Context, name string) error {

	tag, err := db.Exec(ctx, releaseLockSql, name)
	if err != nil {
		return fmt.Errorf("unable to release lock with name %s. %w", name, err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("unable to release lock with name %s, rows affected %d", name, tag.RowsAffected())
	}

	return nil
}

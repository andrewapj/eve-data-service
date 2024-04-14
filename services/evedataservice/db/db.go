package db

import (
	"context"
	"embed"
	"fmt"
	"github.com/andrewapj/arcturus/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"io/fs"
	"strings"
	"sync"
)

//go:embed "schema/*.sql"
var schemaFS embed.FS

var (
	db     *pgxpool.Pool
	dbOnce sync.Once
)

// Connect will connect to the database once. It returns an error if it is unable to connect.
func Connect(ctx context.Context) error {

	var err error
	dbOnce.Do(func() {
		db, err = connect(ctx)
	})

	return err
}

// Close will close the database connection.
func Close() {
	db.Close()
}

// connect performs the database connection. It returns a pointer to the pool or an error.
func connect(ctx context.Context) (*pgxpool.Pool, error) {

	dbPool, err := pgxpool.New(ctx, config.DbUrl())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the database. %w", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to ping the database. %w", err)
	}

	err = generateSchema(ctx, dbPool)
	if err != nil {
		return nil, fmt.Errorf("unable to generate the schema. %w", err)
	}

	return dbPool, nil
}

// generateSchema will look for the required .sql files and execute them to build the database schema.
// It returns an error if it is unable to build it.
func generateSchema(ctx context.Context, db *pgxpool.Pool) error {

	err := fs.WalkDir(schemaFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		fileContent, err := fs.ReadFile(schemaFS, path)
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, string(fileContent))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

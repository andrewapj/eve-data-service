package db

import (
	"context"
	"github.com/andrewapj/arcturus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_connect(t *testing.T) {

	config.SetTestConfig()

	db, err := connect(context.Background())
	require.NoError(t, err, "unexpected error connecting to the database")
	defer db.Close()

	var val int
	err = db.QueryRow(context.Background(), "SELECT 1").Scan(&val)
	require.NoErrorf(t, err, "unexpected error while querying the DB: %v", err)

	assert.Equal(t, 1, val)
}

func Test_generateSchema(t *testing.T) {

	config.SetTestConfig()

	db, err := connect(context.Background())
	require.NoError(t, err, "unexpected error connecting to the database")
	defer db.Close()

	var exists bool
	err = db.QueryRow(context.Background(), `SELECT EXISTS (
   SELECT 1
   FROM   information_schema.tables
   WHERE  table_schema = 'public'
   AND    table_name = 'lock');`).Scan(&exists)
	require.NoError(t, err, "unexpected error while verifying the schema ")

	assert.Truef(t, exists, "expected table 'lock' to exist")
}

func connectToDBAndTruncate(t *testing.T, ctx context.Context) {
	err := Connect(ctx)
	require.NoError(t, err, "unexpected error connecting to the database")

	tables := []string{"expiry", "lock", "status"}

	for _, table := range tables {
		_, err = db.Exec(ctx, "TRUNCATE TABLE "+table)
		require.NoError(t, err, "unexpected error truncating the database")
	}
}

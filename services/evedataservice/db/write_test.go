package db

import (
	"context"
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/config"
	"github.com/andrewapj/arcturus/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_StartTx(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()
	connectToDBAndTruncate(t, ctx)

	tx, err := StartTx(ctx)
	require.NoError(t, err)

	assert.NotNil(t, tx)
}

func Test_CommitTx(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()
	connectToDBAndTruncate(t, ctx)

	tx, err := StartTx(ctx)
	require.NoError(t, err)

	tag, err := tx.Exec(ctx, `INSERT INTO expiry (kind, expires) VALUES ($1, $2)`, "expiry", clock.GetTime())
	require.NoError(t, err)
	assert.Equal(t, int64(1), tag.RowsAffected())

	err = CommitTx(ctx, tx)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, `select count(*) from expiry`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func Test_RollbackTx(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()
	connectToDBAndTruncate(t, ctx)

	tx, err := StartTx(ctx)
	require.NoError(t, err)

	tag, err := tx.Exec(ctx, `INSERT INTO expiry (kind, expires) VALUES ($1, $2)`, "expiry", clock.GetTime())
	require.NoError(t, err)
	assert.Equal(t, int64(1), tag.RowsAffected())

	err = RollbackTx(ctx, tx)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, `select count(*) from expiry`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestInsertBatch(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()
	connectToDBAndTruncate(t, ctx)

	tx, err := StartTx(ctx)
	require.NoError(t, err)

	err = InsertOrUpdateBatch(ctx, tx, "status", []Entity{domain.Status{
		Players:       0,
		ServerVersion: "",
		StartTime:     time.Time{},
		VIP:           false,
	}})
	require.NoError(t, err)

	err = CommitTx(ctx, tx)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, `select count(*) from status`).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

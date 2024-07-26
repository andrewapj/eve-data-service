package db

import (
	"context"
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInsertOrUpdateExpiry_WithMissingExpiry(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()
	connectToDBAndTruncate(t, ctx)

	expected := clock.GetTime()
	err := InsertOrUpdateExpiry(ctx, "item", expected)
	require.NoError(t, err)

	actual, err := FindExpiry(ctx, "item")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestInsertOrUpdateExpiry_WithExistingExpiry(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()
	connectToDBAndTruncate(t, ctx)

	oldExpiry := clock.GetTime().Add(-5 * time.Minute)
	err := InsertOrUpdateExpiry(ctx, "item", oldExpiry)
	require.NoError(t, err)

	expected := clock.GetTime().Add(5 * time.Minute)
	err = InsertOrUpdateExpiry(ctx, "item", expected)
	require.NoError(t, err)

	actual, err := FindExpiry(ctx, "item")
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestFindExpiry(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()
	connectToDBAndTruncate(t, ctx)

	expiry := clock.GetTime()
	err := InsertOrUpdateExpiry(ctx, "item", expiry)
	require.NoError(t, err)

	actual, err := FindExpiry(ctx, "item")
	require.NoError(t, err)
	assert.Equal(t, expiry, actual)
}

func TestFindExpiry_WithNoRowsDefaultTo(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()
	connectToDBAndTruncate(t, ctx)

	actual, err := FindExpiry(ctx, "item")
	require.NoError(t, err)
	assert.Less(t, actual, clock.GetTime())
}

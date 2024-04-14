package db

import (
	"context"
	"github.com/andrewapj/arcturus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
)

const testLockName = "test"

func TestAcquireLock(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()

	connectToDBAndTruncate(t, ctx)

	assert.Truef(t, AcquireLock(ctx, testLockName), "should acquire lock")
}

func TestAcquireLock_FailsWithExistingLock(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()

	connectToDBAndTruncate(t, ctx)

	_, err := db.Exec(ctx, "INSERT INTO lock(name) VALUES ($1)", testLockName)
	require.NoError(t, err)

	assert.Falsef(t, AcquireLock(ctx, testLockName), "should not acquire lock")
}

func TestAcquireLock_Concurrent(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()

	connectToDBAndTruncate(t, ctx)

	var wg sync.WaitGroup
	failures := atomic.Int32{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			if !AcquireLock(ctx, testLockName) {
				failures.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, int32(99), failures.Load())
}

func TestReleaseLock_WithExistingLock(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()

	connectToDBAndTruncate(t, ctx)

	success := AcquireLock(ctx, testLockName)
	require.True(t, success)

	err := ReleaseLock(ctx, testLockName)
	require.NoError(t, err)
}

func TestReleaseLock_WithNoExistingLock(t *testing.T) {

	ctx := context.Background()
	config.SetTestConfig()

	connectToDBAndTruncate(t, ctx)

	err := ReleaseLock(ctx, testLockName)
	require.Error(t, err)
}

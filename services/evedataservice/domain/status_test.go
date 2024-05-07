package domain

import (
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/esi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapStatusFromEsi(t *testing.T) {

	startTime := clock.GetTime()

	esiStatus := &esi.Status{
		Players:       123,
		ServerVersion: "v1",
		StartTime:     startTime,
		VIP:           false,
	}

	status, err := MapStatusFromEsi(esiStatus)
	require.NoError(t, err)

	assert.Equal(t, &Status{
		Players:       123,
		ServerVersion: "v1",
		StartTime:     startTime,
		VIP:           false,
	}, status)
}

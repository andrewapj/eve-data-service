package esi

import (
	"context"
	"github.com/andrewapj/arcturus/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_FetchStatus(t *testing.T) {

	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()

	status, err := client.FetchStatus(context.Background())
	require.NoError(t, err)

	assert.Equal(t, Status{
		Players:       25953,
		ServerVersion: "2539399",
		StartTime:     time.Date(2024, 4, 5, 11, 3, 48, 0, time.UTC).Truncate(time.Second).UTC(),
		VIP:           false,
		baseEsiModel: baseEsiModel{
			expires: esiExpiresTime,
			pages:   1,
		},
	}, status)
}

func BenchmarkClient_FetchStatus(b *testing.B) {

	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()

	for i := 0; i < b.N; i++ {
		_, err := client.FetchStatus(context.Background())
		if err != nil {
			b.Fatal(err.Error())
		}
	}
	b.ReportAllocs()
}

package esi

import (
	"github.com/andrewapj/arcturus/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_mapToIds(t *testing.T) {

	oldTime := clock.GetTime().Add(-time.Duration(5) * time.Minute)
	currentTime := clock.GetTime()

	responses := []*response{
		{body: []byte("[1,2,3]"), expires: currentTime, pages: 2, statusCode: 200},
		{body: []byte("[4,5,6]"), expires: oldTime, pages: 2, statusCode: 200}}

	ids, err := mapToIds(responses)
	require.NoError(t, err)
	require.Equal(t, 6, len(ids.Ids))

	assert.Equal(t, Ids{Ids: []int{1, 2, 3, 4, 5, 6},
		baseEsiModel: baseEsiModel{expires: oldTime, pages: 2}}, ids)
}

func Test_mapToSingle(t *testing.T) {

	currentTime := clock.GetTime()

	responses := []*response{
		{body: []byte(esiStatusResponse), expires: currentTime, pages: 1, statusCode: 200}}

	single, err := mapToSingle[*Status](responses)
	require.NoError(t, err)

	assert.Equal(t, &Status{
		Players: 25953, ServerVersion: "2539399",
		StartTime:    time.Date(2024, 4, 5, 11, 3, 48, 0, time.UTC),
		VIP:          false,
		baseEsiModel: baseEsiModel{expires: currentTime, pages: 1},
	}, single)
}

package clock

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_ParseWithDefault(t *testing.T) {

	loc, err := time.LoadLocation("GMT")
	require.Nil(t, err)
	expected := time.Date(2024, time.March, 31, 22, 7, 36, 0, loc).Truncate(time.Second)

	actual := ParseWithDefault(time.RFC1123, "Sun, 31 Mar 2024 22:07:36 GMT", GetTime()).Truncate(time.Second)

	assert.Equal(t, expected.UTC(), actual.UTC())
}

func Test_ParseWithDefault_WithDefault(t *testing.T) {

	expected := GetTime()
	actual := ParseWithDefault(time.RFC1123, "", expected)

	assert.Equal(t, expected, actual)
}

func Test_GetTime(t *testing.T) {
	assert.Equal(t, time.UTC, GetTime().Location())
}

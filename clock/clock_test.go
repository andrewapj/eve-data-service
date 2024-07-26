package clock

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_GetTime(t *testing.T) {
	assert.Equal(t, time.UTC, GetTime().Location())
}

func TestParseWithDefault(t *testing.T) {
	type args struct {
		layout     string
		input      string
		defaultVal time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "should parse time",
			args: args{
				layout:     time.RFC1123,
				input:      "Sun, 31 Mar 2024 11:05:00 GMT",
				defaultVal: GetTime().Truncate(time.Second),
			},
			want: time.Date(2024, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second).UTC(),
		},
		{
			name: "should not parse invalid time and return a default",
			args: args{
				layout:     time.RFC1123,
				input:      "",
				defaultVal: time.Date(2024, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second),
			},
			want: time.Date(2024, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second).UTC(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ParseWithDefault(tt.args.layout, tt.args.input, tt.args.defaultVal), "ParseWithDefault(%v, %v, %v)", tt.args.layout, tt.args.input, tt.args.defaultVal)
		})
	}
}

func TestFindEarliestTime(t *testing.T) {

	_, err := FindEarliestTime(nil)
	require.Error(t, err)

	_, err = FindEarliestTime([]time.Time{})

	actual, err := FindEarliestTime([]time.Time{
		time.Date(2024, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second),
	})
	require.NoError(t, err)
	assert.Equal(t, time.Date(2024, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second), actual)

	actual, err = FindEarliestTime([]time.Time{
		time.Date(2024, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second),
		time.Date(2023, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second),
		time.Date(2025, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second),
	})
	require.NoError(t, err)
	assert.Equal(t, time.Date(2023, time.March, 31, 11, 5, 0, 0, time.UTC).Truncate(time.Second), actual)
}

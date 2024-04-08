package clock

import (
	"fmt"
	"time"
)

// GetTime gets the current time in UTC.
func GetTime() time.Time {
	return time.Now().UTC()
}

// ParseWithDefault parses a time string and returns a default if it can not be parsed.
func ParseWithDefault(layout string, input string, defaultVal time.Time) time.Time {

	val, err := time.ParseInLocation(layout, input, time.UTC)
	if err != nil {
		return defaultVal.UTC()
	}
	return val.UTC()
}

// FindEarliestTime finds the earliest time from a slice of times.
func FindEarliestTime(t []time.Time) (time.Time, error) {
	if len(t) == 0 {
		return time.Time{}, fmt.Errorf("cannot find earliest time from an empty slice")
	}

	earliest := t[0]

	for _, t := range t[1:] {
		if t.Before(earliest) {
			earliest = t
		}
	}

	return earliest, nil
}

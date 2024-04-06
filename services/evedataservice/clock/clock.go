package clock

import "time"

// ParseWithDefault parses a time string and returns a default if it can not be parsed.
func ParseWithDefault(layout string, input string, defaultVal time.Time) time.Time {

	val, err := time.ParseInLocation(layout, input, time.UTC)
	if err != nil {
		return defaultVal.UTC()
	}
	return val.UTC()
}

// GetTime gets the current time in UTC.
func GetTime() time.Time {
	return time.Now().UTC()
}

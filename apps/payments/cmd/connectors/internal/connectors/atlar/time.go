package atlar

import "time"

func ParseAtlarTimestamp(value string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, value)
}

func ParseAtlarDate(value string) (time.Time, error) {
	return time.Parse(time.DateOnly, value)
}

package atlar

import "time"

func ParseAtlarTimestamp(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.999999999Z", value)
}

func ParseAtlarDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

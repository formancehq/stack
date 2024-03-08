package storage

import "regexp"

var (
	metadataRegex = regexp.MustCompile("metadata\\[(.+)\\]")
)

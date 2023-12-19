package client

import (
	"errors"
	"regexp"
)

const (
	atlarMetadataSpecNamespace = "com.atlar.spec/"
)

func ExtractNamespacedMetadataIgnoreEmpty(metadata map[string]string, key string) *string {
	value := metadata[atlarMetadataSpecNamespace+key]
	return &value
}

type IdentifierData struct {
	Market string
	Type   string
	Number string
}

var identifierMetadataRegex = regexp.MustCompile(`^com\.atlar\.spec/identifier/([^/]+)/([^/]+)$`)

func metadataToIdentifierData(key, value string) (*IdentifierData, error) {
	// Find matches in the input string
	matches := identifierMetadataRegex.FindStringSubmatch(key)
	if matches == nil {
		return nil, errors.New("input does not match the expected format")
	}

	// Extract values from the matched groups
	return &IdentifierData{
		Market: matches[1],
		Type:   matches[2],
		Number: value,
	}, nil
}

package internal

import (
	"os"
)

func getOpenSearchUrl() string {
	if fromEnv := os.Getenv("OPENSEARCH_URL"); fromEnv != "" {
		return fromEnv
	}
	return "localhost:9200"
}

func getTemporalAddress() string {
	if fromEnv := os.Getenv("TEMPORAL_ADDRESS"); fromEnv != "" {
		return fromEnv
	}
	return "localhost:7233"
}

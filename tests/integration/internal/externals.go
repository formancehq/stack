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

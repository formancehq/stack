package internal

import (
	"net/url"
	"os"

	. "github.com/onsi/gomega"
)

func getOpenSearchHost() string {
	openSearchUrl := getOpenSearchUrl()
	url, err := url.Parse(openSearchUrl)
	Expect(err).To(BeNil())
	return url.Host
}

func getOpenSearchUrl() string {
	if fromEnv := os.Getenv("OPENSEARCH_URL"); fromEnv != "" {
		return fromEnv
	}
	return "localhost:9200"
}

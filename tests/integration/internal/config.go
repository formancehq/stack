package internal

import "os"

func GetPostgresDSNString() string {
	if fromEnv := os.Getenv("POSTGRES_DSN"); fromEnv != "" {
		return fromEnv
	}
	return "postgres://formance:formance@localhost:5432/formance?sslmode=disable"
}

func GetNatsAddress() string {
	return "localhost:4222" // TODO: Make configurable
}

func GetOpenSearchUrl() string {
	if fromEnv := os.Getenv("OPENSEARCH_URL"); fromEnv != "" {
		return fromEnv
	}
	return "localhost:9200"
}

func GetTemporalAddress() string {
	if fromEnv := os.Getenv("TEMPORAL_ADDRESS"); fromEnv != "" {
		return fromEnv
	}
	return "localhost:7233"
}

func GetDockerEndpoint() string {
	host := os.Getenv("DOCKER_HOSTNAME")
	if host == "" {
		host = "host.docker.internal"
	}
	return host
}

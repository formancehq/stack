package v1beta3

import (
	"time"
)

type LockingStrategyRedisConfig struct {
	Uri string `json:"uri,omitempty"`
	// +optional
	TLS bool `json:"tls"`
	// +optional
	InsecureTLS bool `json:"insecure,omitempty"`
	// +optional
	Duration time.Duration `json:"duration,omitempty"`
	// +optional
	Retry time.Duration `json:"retry,omitempty"`
}

// +kubebuilder:object:generate=true
type LedgerSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	AllowPastTimestamps bool `json:"allowPastTimestamps"`
}

package v1beta3

import (
	"github.com/formancehq/operator/v2/api/formance.com/v1beta1"
	"time"
)

type LockingStrategyRedisConfig struct {
	// +optional
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

type LockingStrategy struct {
	// +kubebuilder:Enum:={memory,redis}
	// +kubebuilder:default:=memory
	// +optional
	Strategy string `json:"strategy,omitempty"`
	// +optional
	Redis LockingStrategyRedisConfig `json:"redis"`
}

// +kubebuilder:object:generate=true
type LedgerSpec struct {
	CommonServiceProperties `json:",inline"`
	Postgres                v1beta1.DatabaseConfigurationSpec `json:"postgres"`
	// +optional
	Locking LockingStrategy `json:"locking"`
	// +optional
	AllowPastTimestamps bool `json:"allowPastTimestamps"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
	// +optional
	DeploymentStrategy v1beta1.DeploymentStrategy `json:"deploymentStrategy,omitempty"`
	// +optional
	Disabled *bool `json:"disabled,omitempty"`
}

type ServiceSpec struct {
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

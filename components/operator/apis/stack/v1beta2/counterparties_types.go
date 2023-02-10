package v1beta2

import (
	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
)

// +kubebuilder:object:generate=true
type CounterpartiesSpec struct {
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// +optional
	Debug bool `json:"debug,omitempty"`
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +optional
	Postgres componentsv1beta2.PostgresConfig `json:"postgres"`
}

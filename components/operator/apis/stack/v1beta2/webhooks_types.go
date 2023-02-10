package v1beta2

import (
	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
)

// +kubebuilder:object:generate=true
type WebhooksSpec struct {
	// +optional
	Debug bool `json:"debug,omitempty"`
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
	// +optional
	Postgres componentsv1beta2.PostgresConfig `json:"postgres"`
}

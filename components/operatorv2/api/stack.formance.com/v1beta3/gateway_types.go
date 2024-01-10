package v1beta3

import "github.com/formancehq/operator/v2/api/formance.com/v1beta1"

// +kubebuilder:object:generate=true
type GatewaySpec struct {
	// +optional
	ResourceProperties *v1beta1.ResourceProperties `json:"resourceProperties,omitempty"`
	// +optional
	Annotations AnnotationsServicesSpec `json:"annotations,omitempty"`
	// +optional
	Fallback *string `json:"fallback,omitempty"`

	// +optional
	EnableAuditPlugin *bool `json:"enableAuditPlugin,omitempty"`

	// +optional
	LivenessEndpoint string `json:"livenessEndpoint,omitempty"`

	// +optional
	EnableScopes *bool `json:"enableScopes,omitempty"`
}

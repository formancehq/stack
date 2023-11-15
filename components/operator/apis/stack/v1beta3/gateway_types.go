package v1beta3

// +kubebuilder:object:generate=true
type GatewaySpec struct {
	// +optional
	ResourceProperties *ResourceProperties `json:"resourceProperties,omitempty"`
	// +optional
	Annotations AnnotationsServicesSpec `json:"annotations,omitempty"`
	// +optional
	Fallback *string `json:"fallback,omitempty"`

	// +optional
	EnableAuditPlugin *bool `json:"enableAuditPlugin,omitempty"`
}

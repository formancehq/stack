package v1beta3

// +kubebuilder:object:generate=true
type GatewaySpec struct {
	// +optional
	ResourceProperties *ResourceProperties `json:"resourceProperties,omitempty"`
	// +optional
	Annotations AnnotationsServicesSpec `json:"annotations,omitempty"`
}

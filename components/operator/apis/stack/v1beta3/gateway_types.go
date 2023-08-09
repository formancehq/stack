package v1beta3

// +kubebuilder:object:generate=true
type GatewaySpec struct {
	// +optional
	Annotations AnnotationsServicesSpec `json:"annotations,omitempty"`
}

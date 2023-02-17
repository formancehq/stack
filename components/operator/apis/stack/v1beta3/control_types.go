package v1beta3

// +kubebuilder:object:generate=true
type ControlSpec struct {
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +kubebuilder:validation:Optional
	// +nullable
	Ingress *IngressConfig `json:"ingress"`
}

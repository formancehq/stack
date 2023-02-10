package v1beta3

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	DevProperties `json:",inline"`
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
}

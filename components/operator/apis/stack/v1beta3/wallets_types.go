package v1beta3

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	DevProperties `json:",inline"`
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +kubebuilder:validation:Optional
	// +nullable
	Ingress *IngressConfig `json:"ingress"`
}

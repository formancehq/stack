package v1beta3

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	// +optional
	DevProperties `json:",inline"`
	Annotations   AnnotationsServicesSpec `json:"annotations,omitempty"`
}

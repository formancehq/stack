package v1beta3

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	// +optional
	DevProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

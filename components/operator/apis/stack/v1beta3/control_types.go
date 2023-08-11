package v1beta3

// +kubebuilder:object:generate=true
type ControlSpec struct {
	// +optional
	DevProperties `json:",inline"`
	Annotations   AnnotationsServicesSpec `json:"annotations,omitempty"`
}

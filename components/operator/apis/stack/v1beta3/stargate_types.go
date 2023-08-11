package v1beta3

type StargateSpec struct {
	// +optional
	DevProperties `json:",inline"`
	Annotations   AnnotationsServicesSpec `json:"annotations,omitempty"`
}

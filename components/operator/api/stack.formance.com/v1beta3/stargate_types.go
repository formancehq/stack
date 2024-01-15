package v1beta3

type StargateSpec struct {
	// +optional
	DevProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

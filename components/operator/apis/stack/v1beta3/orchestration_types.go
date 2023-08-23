package v1beta3

// +kubebuilder:object:generate=true
type OrchestrationSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	DevProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

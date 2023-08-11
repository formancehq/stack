package v1beta3

// +kubebuilder:object:generate=true
type OrchestrationSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	DevProperties `json:",inline"`
	Annotations   AnnotationsServicesSpec `json:"annotations,omitempty"`
}

package v1beta3

// +kubebuilder:object:generate=true
type OrchestrationSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	Annotations AnnotationsServicesSpec `json:"service"`
}

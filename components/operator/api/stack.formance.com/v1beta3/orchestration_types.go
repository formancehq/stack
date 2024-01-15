package v1beta3

// +kubebuilder:object:generate=true
type OrchestrationSpec struct {
	CommonServiceProperties `json:",inline"`
	Postgres                PostgresConfig `json:"postgres"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

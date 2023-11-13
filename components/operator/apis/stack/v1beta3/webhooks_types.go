package v1beta3

// +kubebuilder:object:generate=true
type WebhooksSpec struct {
	CommonServiceProperties `json:",inline"`
	Postgres                PostgresConfig `json:"postgres"`
	// +optional
	DevProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

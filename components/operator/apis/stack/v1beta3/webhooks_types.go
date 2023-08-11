package v1beta3

// +kubebuilder:object:generate=true
type WebhooksSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	DevProperties `json:",inline"`
	Annotations   AnnotationsServicesSpec `json:"annotations,omitempty"`
}

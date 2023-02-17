package v1beta3

// +kubebuilder:object:generate=true
type OrchestrationSpec struct {
	DevProperties `json:",inline"`
	// +optional
	Postgres PostgresConfig `json:"postgres"`
}

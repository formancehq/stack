package v1beta3

// +kubebuilder:object:generate=true
type OrchestrationSpec struct {
	DevProperties `json:",inline"`
	Postgres      PostgresConfig `json:"postgres"`
}

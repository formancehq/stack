package v1beta3

// +kubebuilder:object:generate=true
type WebhooksSpec struct {
	DevProperties `json:",inline"`
	Postgres      PostgresConfig `json:"postgres"`
}

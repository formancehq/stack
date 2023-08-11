package v1beta3

type AuthSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	StaticClients []StaticClient `json:"staticClients,omitempty"`

	// +optional
	DevProperties `json:",inline"`
	Annotations   AnnotationsServicesSpec `json:"annotations,omitempty"`
}

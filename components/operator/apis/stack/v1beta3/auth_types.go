package v1beta3

type AuthSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	StaticClients []StaticClient `json:"staticClients,omitempty"`

	// +optional
	DevProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

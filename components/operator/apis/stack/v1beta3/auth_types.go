package v1beta3

type AuthSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	StaticClients []StaticClient `json:"staticClients,omitempty"`

	// +optional
	Annotations AnnotationsServicesSpec `json:"annotations,omitempty"`
}

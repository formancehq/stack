package v1beta3

type AuthSpec struct {
	Postgres PostgresConfig `json:"postgres"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
	// +optional
	StaticClients []StaticClient `json:"staticClients"`
}

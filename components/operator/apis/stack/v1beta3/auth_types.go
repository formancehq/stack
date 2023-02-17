package v1beta3

type AuthSpec struct {
	// +kubebuilder:validation:Required
	// +nullable
	Postgres PostgresConfig `json:"postgres"`
	// +kubebuilder:validation:Optional
	// +nullable
	Ingress *IngressConfig `json:"ingress,omitempty"`
	// +kubebuilder:validation:Optional
	// +nullable
	StaticClients []StaticClient `json:"staticClients,omitempty"`
}

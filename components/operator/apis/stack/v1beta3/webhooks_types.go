package v1beta3

// +kubebuilder:object:generate=true
type WebhooksSpec struct {
	// +optional
	Debug bool `json:"debug,omitempty"`
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
	// +optional
	Postgres PostgresConfig `json:"postgres"`
}

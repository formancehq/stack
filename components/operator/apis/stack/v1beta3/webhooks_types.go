package v1beta3

// +kubebuilder:object:generate=true
type WebhooksSpec struct {
	// +optional
	Debug bool `json:"debug,omitempty"`
	// +optional
	Postgres PostgresConfig `json:"postgres"`
}

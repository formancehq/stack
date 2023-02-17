package v1beta3

// +kubebuilder:object:generate=true
type PaymentsSpec struct {
	EncryptionKey string `json:"encryptionKey"`
	// +optional
	Postgres PostgresConfig `json:"postgres"`
}

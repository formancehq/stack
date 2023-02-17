package v1beta3

// +kubebuilder:object:generate=true
type PaymentsSpec struct {
	EncryptionKey string         `json:"encryptionKey"`
	Postgres      PostgresConfig `json:"postgres"`
}

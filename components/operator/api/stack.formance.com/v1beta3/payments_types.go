package v1beta3

// +kubebuilder:object:generate=true
type PaymentsSpec struct {
	CommonServiceProperties `json:",inline"`
	EncryptionKey           string         `json:"encryptionKey"`
	Postgres                PostgresConfig `json:"postgres"`

	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

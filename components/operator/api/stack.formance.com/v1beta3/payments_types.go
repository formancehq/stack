package v1beta3

import "github.com/formancehq/operator/api/formance.com/v1beta1"

// +kubebuilder:object:generate=true
type PaymentsSpec struct {
	CommonServiceProperties `json:",inline"`
	EncryptionKey           string                            `json:"encryptionKey"`
	Postgres                v1beta1.DatabaseConfigurationSpec `json:"postgres"`

	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

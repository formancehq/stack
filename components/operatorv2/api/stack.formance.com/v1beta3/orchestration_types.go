package v1beta3

import "github.com/formancehq/operator/v2/api/formance.com/v1beta1"

// +kubebuilder:object:generate=true
type OrchestrationSpec struct {
	CommonServiceProperties `json:",inline"`
	Postgres                v1beta1.DatabaseConfigurationSpec `json:"postgres"`
	// +optional
	ResourceProperties *v1beta1.ResourceProperties `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec     `json:"annotations,omitempty"`
}

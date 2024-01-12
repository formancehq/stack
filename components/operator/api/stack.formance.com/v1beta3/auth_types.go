package v1beta3

import "github.com/formancehq/operator/api/formance.com/v1beta1"

type AuthSpec struct {
	Postgres DatabaseConfigurationSpec `json:"postgres"`
	// +optional
	StaticClients []StaticClient `json:"staticClients,omitempty"`

	// +optional
	v1beta1.DevProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

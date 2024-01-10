package v1beta3

import "github.com/formancehq/operator/v2/api/formance.com/v1beta1"

type AuthSpec struct {
	Postgres v1beta1.DatabaseConfigurationSpec `json:"postgres"`
	// +optional
	StaticClients []StaticClient `json:"staticClients,omitempty"`

	// +optional
	v1beta1.DevProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

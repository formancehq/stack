package v1beta3

import "github.com/formancehq/operator/v2/api/formance.com/v1beta1"

type StargateSpec struct {
	// +optional
	*v1beta1.DevProperties `json:",inline"`
	// +optional
	ResourceProperties *ResourceProperties     `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec `json:"annotations,omitempty"`
}

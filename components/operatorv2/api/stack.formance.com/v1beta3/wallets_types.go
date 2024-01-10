package v1beta3

import "github.com/formancehq/operator/v2/api/formance.com/v1beta1"

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	CommonServiceProperties `json:",inline"`
	// +optional
	ResourceProperties *v1beta1.ResourceProperties `json:"resourceProperties,omitempty"`
	Annotations        AnnotationsServicesSpec     `json:"annotations,omitempty"`
}

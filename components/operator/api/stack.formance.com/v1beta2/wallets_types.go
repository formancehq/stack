package v1beta2

import (
	componentsv1beta2 "github.com/formancehq/operator/api/stack.formance.com/components/v1beta2"
)

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	componentsv1beta2.DevProperties `json:",inline"`
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
}

package v1beta2

import (
	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
)

// +kubebuilder:object:generate=true
type SearchSpec struct {
	ElasticSearchConfig componentsv1beta2.ElasticSearchConfig `json:"elasticSearch"`

	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`

	//+optional
	Ingress *IngressConfig `json:"ingress"`

	// +optional
	Batching componentsv1beta2.Batching `json:"batching"`
}

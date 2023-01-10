package v1beta2

import (
	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
	"github.com/formancehq/operator/pkg/apis/v1beta2"
	"github.com/formancehq/operator/pkg/typeutils"
	"k8s.io/apimachinery/pkg/util/validation/field"
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

func (in *SearchSpec) Validate() field.ErrorList {
	if in == nil {
		return nil
	}
	return typeutils.Map(in.ElasticSearchConfig.Validate(), v1beta2.AddPrefixToFieldError("elasticSearch"))
}

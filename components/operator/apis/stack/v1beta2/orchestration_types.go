package v1beta2

import (
	"github.com/formancehq/operator/pkg/apis/v1beta2"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +kubebuilder:object:generate=true
type OrchestrationSpec struct {
	v1beta2.DevProperties `json:",inline"`
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
	// +optional
	Postgres v1beta2.PostgresConfig `json:"postgres"`
}

func (in *OrchestrationSpec) Validate() field.ErrorList {
	return field.ErrorList{}
}

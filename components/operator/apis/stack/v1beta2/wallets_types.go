package v1beta2

import (
	"github.com/formancehq/operator/pkg/apis/v1beta2"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// +kubebuilder:object:generate=true
type WalletsSpec struct {
	v1beta2.DevProperties `json:",inline"`
	// +optional
	Scaling ScalingSpec `json:"scaling,omitempty"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
}

func (in *WalletsSpec) Validate() field.ErrorList {
	return field.ErrorList{}
}

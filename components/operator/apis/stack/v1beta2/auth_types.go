package v1beta2

import (
	authcomponentsv1beta2 "github.com/formancehq/operator/apis/auth.components/v1beta2"
	"github.com/formancehq/operator/pkg/apis/v1beta2"
	"github.com/formancehq/operator/pkg/typeutils"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type AuthSpec struct {
	Postgres v1beta2.PostgresConfig `json:"postgres"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
	// +optional
	StaticClients []authcomponentsv1beta2.StaticClient `json:"staticClients"`
}

func (in *AuthSpec) Validate() field.ErrorList {
	if in == nil {
		return field.ErrorList{}
	}
	return typeutils.MergeAll(
		typeutils.Map(in.Postgres.Validate(), v1beta2.AddPrefixToFieldError("postgres.")),
	)
}

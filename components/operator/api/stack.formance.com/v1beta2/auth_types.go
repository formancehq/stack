package v1beta2

import (
	authcomponentsv1beta2 "github.com/formancehq/operator/api/stack.formance.com/auth.components/v1beta2"
	componentsv1beta2 "github.com/formancehq/operator/api/stack.formance.com/components/v1beta2"
)

type AuthSpec struct {
	Postgres componentsv1beta2.PostgresConfig `json:"postgres"`
	// +optional
	Ingress *IngressConfig `json:"ingress"`
	// +optional
	StaticClients []authcomponentsv1beta2.StaticClient `json:"staticClients"`
}

func (in AuthSpec) NeedAuthMiddleware() bool {
	return false
}

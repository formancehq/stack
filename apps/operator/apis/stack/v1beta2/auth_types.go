package v1beta2

import (
	"fmt"
	"sort"
	"strings"

	authcomponentsv1beta2 "github.com/formancehq/operator/apis/auth.components/v1beta2"
	componentsv1beta2 "github.com/formancehq/operator/apis/components/v1beta2"
	"github.com/formancehq/operator/internal/collectionutils"
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

func (in AuthSpec) Spec(stack *Stack, configuration ConfigurationSpec) any {
	stackStaticClients := collectionutils.SliceFromMap(stack.Status.StaticAuthClients)
	sort.SliceStable(
		stackStaticClients,
		func(i, j int) bool {
			return strings.Compare(stackStaticClients[i].ID, stackStaticClients[j].ID) < 0
		},
	)
	staticClients := append(configuration.Services.Auth.StaticClients, stackStaticClients...)
	staticClients = append(staticClients, stack.Spec.Auth.StaticClients...)
	return componentsv1beta2.AuthSpec{
		Postgres: componentsv1beta2.PostgresConfigCreateDatabase{
			CreateDatabase: true,
			PostgresConfigWithDatabase: componentsv1beta2.PostgresConfigWithDatabase{
				PostgresConfig: configuration.Services.Auth.Postgres,
				Database:       fmt.Sprintf("%s-auth", stack.Name),
			},
		},
		BaseURL:             fmt.Sprintf("%s://%s/api/auth", stack.Spec.Scheme, stack.Spec.Host),
		DelegatedOIDCServer: stack.Spec.Auth.DelegatedOIDCServer,
		StaticClients:       staticClients,
	}
}

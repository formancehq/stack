package gateways

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func createConfigMap(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs []*v1beta1.GatewayHTTPAPI, auth *v1beta1.Auth, auditTopic *v1beta1.BrokerTopic) (*v1.ConfigMap, error) {

	caddyfile, err := CreateCaddyfile(ctx, stack, gateway, httpAPIs, auth, auditTopic)
	if err != nil {
		return nil, err
	}

	caddyfileConfigMap, _, err := core.CreateOrUpdate[*v1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "gateway",
	},
		func(t *v1.ConfigMap) error {
			t.Data = map[string]string{
				"Caddyfile": caddyfile,
			}

			return nil
		},
		core.WithController[*v1.ConfigMap](ctx.GetScheme(), gateway),
	)

	return caddyfileConfigMap, err
}

package gateways

import (
	"github.com/formancehq/go-libs/collectionutils"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/caddy"
)

func CreateCaddyfile(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs []*v1beta1.GatewayHTTPAPI, broker *v1beta1.Broker) (string, error) {

	data := map[string]any{
		"Services": collectionutils.Map(httpAPIs, func(from *v1beta1.GatewayHTTPAPI) v1beta1.GatewayHTTPAPISpec {
			return from.Spec
		}),
		"Platform": ctx.GetPlatform(),
		"Debug":    stack.Spec.Debug,
		"Port":     8080,
		"Gateway": map[string]any{
			"Version": gateway.Spec.Version,
		},
	}

	// TODO(gfyrag): Check if search is enabled
	if stack.Spec.EnableAudit && broker != nil {
		data["EnableAudit"] = true
		data["Broker"] = broker.Status.URI.Scheme
	}

	return caddy.ComputeCaddyfile(ctx, stack, Caddyfile, data)
}

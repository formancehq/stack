package gateways

import (
	"fmt"

	"github.com/formancehq/operator/internal/resources/settings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
)

func CreateCaddyfile(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs []*v1beta1.GatewayHTTPAPI, auth *v1beta1.Auth, auditTopic *v1beta1.BrokerTopic) (string, error) {

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
	if auth != nil {
		data["Auth"] = map[string]any{
			"Issuer":       fmt.Sprintf("%s/api/auth", URL(gateway)),
			"EnableScopes": auth.Spec.EnableScopes,
		}
	}

	if stack.Spec.EnableAudit && auditTopic != nil {
		data["EnableAudit"] = true
		data["Broker"] = auditTopic.Status.URI.Scheme
	}

	return settings.ComputeCaddyfile(ctx, stack, Caddyfile, data)
}

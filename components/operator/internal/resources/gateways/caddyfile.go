package gateways

import (
	"fmt"
	"sort"

	"github.com/formancehq/operator/internal/resources/settings"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
)

func CreateCaddyfile(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, httpAPIs []*v1beta1.HTTPAPI, auth *v1beta1.Auth, auditTopic *v1beta1.BrokerTopic) (string, error) {

	sort.Slice(httpAPIs, func(i, j int) bool {
		return httpAPIs[i].Spec.Name < httpAPIs[j].Spec.Name
	})

	data := map[string]any{
		"Services": collectionutils.Map(httpAPIs, func(from *v1beta1.HTTPAPI) v1beta1.HTTPAPISpec {
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
		data["Broker"] = func() string {
			if auditTopic.Status.Configuration.Kafka != nil {
				return "kafka"
			}
			if auditTopic.Status.Configuration.Nats != nil {
				return "nats"
			}
			return ""
		}()
	}

	return settings.ComputeCaddyfile(ctx, stack, Caddyfile, data)
}

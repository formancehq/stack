package gateways

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	v1 "k8s.io/api/core/v1"
)

func createDeployment(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, caddyfileConfigMap *v1.ConfigMap,
	auditTopic *v1beta1.BrokerTopic, version string) error {

	env := GetEnvVars(gateway)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseKind(ctx, gateway))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, core.GetDevEnvVars(stack, gateway)...)

	if stack.Spec.EnableAudit && auditTopic != nil {
		env = append(env, settings.GetBrokerEnvVars(auditTopic.Status.URI, stack.Name, "gateway")...)
	}

	image, err := registries.GetImage(ctx, stack, "gateway", version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, stack, gateway, "gateway",
		settings.ConfigureCaddy(caddyfileConfigMap, image, env),
		deployments.WithMatchingLabels("gateway"),
	)
	return err
}

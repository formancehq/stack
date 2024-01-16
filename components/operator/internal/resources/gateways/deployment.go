package gateways

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/brokerconfigurations"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/opentelemetryconfigurations"
	"github.com/formancehq/operator/internal/resources/registries"
	v1 "k8s.io/api/core/v1"
)

func createDeployment(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, caddyfileConfigMap *v1.ConfigMap, auditTopic *v1beta1.BrokerTopic) error {

	env := make([]v1.EnvVar, 0)
	otlpEnv, err := opentelemetryconfigurations.EnvVarsIfEnabled(ctx, stack.Name, core.GetModuleName(gateway))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, core.GetDevEnvVars(stack, gateway)...)

	if stack.Spec.EnableAudit && auditTopic != nil {
		env = append(env,
			brokerconfigurations.BrokerEnvVars(*auditTopic.Status.Configuration, stack.Name, "gateway")...,
		)
	}

	image, err := registries.GetImage(ctx, stack, "gateway", gateway.Spec.Version)
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, gateway, "gateway",
		core.ConfigureCaddy(caddyfileConfigMap, image, env, gateway.Spec.ResourceRequirements),
		deployments.WithMatchingLabels("gateway"),
	)
	return err
}

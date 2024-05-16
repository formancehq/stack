package gateways

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/applications"
	"github.com/formancehq/operator/internal/resources/brokers"
	"github.com/formancehq/operator/internal/resources/caddy"
	"github.com/formancehq/operator/internal/resources/registries"
	v1 "k8s.io/api/core/v1"
)

func createDeployment(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, caddyfileConfigMap *v1.ConfigMap,
	broker *v1beta1.Broker, version string) error {

	env := GetEnvVars(gateway)
	env = append(env, core.GetDevEnvVars(stack, gateway)...)

	if stack.Spec.EnableAudit && broker != nil {
		brokerEnvVar, err := brokers.GetBrokerEnvVars(ctx, broker.Status.URI, stack.Name, "gateway")
		if err != nil {
			return err
		}

		env = append(env, brokerEnvVar...)
	}

	image, err := registries.GetImage(ctx, stack, "gateway", version)
	if err != nil {
		return err
	}

	caddyTpl, err := caddy.DeploymentTemplate(ctx, stack, gateway, caddyfileConfigMap, image, env)
	if err != nil {
		return err
	}

	caddyTpl.Name = "gateway"
	return applications.
		New(gateway, caddyTpl).
		IsEE().
		Install(ctx)
}

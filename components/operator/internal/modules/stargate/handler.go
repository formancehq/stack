package stargate

import (
	"strconv"
	"strings"

	"github.com/formancehq/operator/internal/modules/gateway"

	"github.com/formancehq/operator/internal/modules"
)

type module struct{}

func (s module) DependsOn() []modules.Module {
	// todo: need a strict requirement
	return []modules.Module{
		gateway.Module,
	}
}

func (s module) Name() string {
	return "stargate"
}

func (s module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ReconciliationConfig) modules.Services {
				if ctx.Stack.Spec.Stargate == nil {
					return modules.Services{}
				}

				return modules.Services{
					{
						ListenEnvVar:       "BIND",
						HasVersionEndpoint: true,
						Liveness:           modules.LivenessDefault,
						Annotations:        ctx.Configuration.Spec.Services.Stargate.Annotations.Service,
						Container: func(resolveContext modules.ContainerResolutionConfiguration) modules.Container {
							return modules.Container{
								Env:   stargateClientEnvVars(resolveContext),
								Image: modules.GetImage("stargate", resolveContext.Versions.Spec.Stargate),
								Resources: modules.GetResourcesWithDefault(
									resolveContext.Configuration.Spec.Services.Stargate.ResourceProperties,
									modules.ResourceSizeSmall(),
								),
							}
						},
					},
				}
			},
		},
	}
}

var Module = &module{}

var _ modules.Module = Module
var _ modules.DependsOnAwareModule = Module

func stargateClientEnvVars(resolveContext modules.ContainerResolutionConfiguration) modules.ContainerEnv {
	l := strings.Split(resolveContext.Stack.ObjectMeta.Name, "-")
	organizationID := l[0]
	stackID := l[1]

	return modules.NewEnv().Append(
		modules.Env("ORGANIZATION_ID", organizationID),
		modules.Env("STACK_ID", stackID),
		modules.Env("STARGATE_SERVER_URL", resolveContext.Stack.Spec.Stargate.StargateServerURL),
		// TODO: The ports of the service must be made available to this part of the code
		modules.Env("GATEWAY_URL", "http://gateway:"+strconv.Itoa(int(resolveContext.RegisteredModules["gateway"].Services["gateway"].Port))),
		modules.Env("STARGATE_AUTH_CLIENT_ID", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.ClientID),
		modules.Env("STARGATE_AUTH_CLIENT_SECRET", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.ClientSecret),
		modules.Env("STARGATE_AUTH_ISSUER_URL", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.Issuer),
	)
}

func init() {
	modules.Register(Module)
}

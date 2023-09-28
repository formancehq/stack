package stargate

import (
	"strconv"
	"strings"

	"github.com/formancehq/operator/internal/modules"
)

type module struct{}

func (s module) Versions() map[string]modules.Version {
	return map[string]modules.Version{
		"v0.0.0": {
			Services: func(ctx modules.ModuleContext) modules.Services {
				if ctx.Stack.Spec.Stargate == nil {
					return modules.Services{}
				}

				return modules.Services{
					{
						ListenEnvVar:       "BIND",
						HasVersionEndpoint: true,
						ExposeHTTP:         true,
						Liveness:           modules.LivenessDefault,
						Annotations:        ctx.Configuration.Spec.Services.Stargate.Annotations.Service,
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
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

var _ modules.Module = (*module)(nil)

func stargateClientEnvVars(resolveContext modules.ContainerResolutionContext) modules.ContainerEnv {
	l := strings.Split(resolveContext.Stack.ObjectMeta.Name, "-")
	organizationID := l[0]
	stackID := l[1]

	return modules.NewEnv().Append(
		modules.Env("ORGANIZATION_ID", organizationID),
		modules.Env("STACK_ID", stackID),
		modules.Env("STARGATE_SERVER_URL", resolveContext.Stack.Spec.Stargate.StargateServerURL),
		// TODO: The ports of the service must be made available to this part of the code
		modules.Env("GATEWAY_URL", "http://gateway:"+strconv.Itoa(int(resolveContext.RegisteredModules["gateway"].Module.Versions()["v0.0.0"].Services(resolveContext.ModuleContext)[0].Port))),
		modules.Env("STARGATE_AUTH_CLIENT_ID", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.ClientID),
		modules.Env("STARGATE_AUTH_CLIENT_SECRET", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.ClientSecret),
		modules.Env("STARGATE_AUTH_ISSUER_URL", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.Issuer),
	)
}

func init() {
	modules.Register("stargate", &module{})
}

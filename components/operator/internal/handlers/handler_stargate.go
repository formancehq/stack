package handlers

import (
	"strconv"
	"strings"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

func init() {
	stargateClientEnvVars := func(resolveContext modules.ContainerResolutionContext) modules.ContainerEnv {
		l := strings.Split(resolveContext.Stack.ObjectMeta.Name, "-")
		organizationID := l[0]
		stackID := l[1]

		return modules.NewEnv().Append(
			modules.Env("ORGANIZATION_ID", organizationID),
			modules.Env("STACK_ID", stackID),
			modules.Env("STARGATE_SERVER_URL", resolveContext.Stack.Spec.Stargate.StargateServerURL),
			// TODO: The ports of the service must be made available to this part of the code
			modules.Env("GATEWAY_URL", "http://gateway:"+strconv.Itoa(int(resolveContext.RegisteredModules["gateway"].Module.Versions["v0.0.0"].Services(resolveContext.ModuleContext)[0].Port))),
			modules.Env("STARGATE_AUTH_CLIENT_ID", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.ClientID),
			modules.Env("STARGATE_AUTH_CLIENT_SECRET", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.ClientSecret),
			modules.Env("STARGATE_AUTH_ISSUER_URL", resolveContext.Stack.Spec.Auth.DelegatedOIDCServer.Issuer),
		)
	}
	modules.Register("stargate", modules.Module{
		Versions: map[string]modules.Version{
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
							AuthConfiguration: func(resolveContext modules.ModuleContext) stackv1beta3.ClientConfiguration {
								return stackv1beta3.NewClientConfiguration()
							},
							Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
								return modules.Container{
									Env:   stargateClientEnvVars(resolveContext),
									Image: modules.GetImage("stargate", resolveContext.Versions.Spec.Stargate),
									Resources: getResourcesWithDefault(
										resolveContext.Configuration.Spec.Services.Stargate.ResourceProperties,
										modules.ResourceSizeSmall(),
									),
								}
							},
						},
					}
				},
			},
		},
	})
}

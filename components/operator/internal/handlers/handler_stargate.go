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
			modules.Env("GATEWAY_URL", "http://gateway:"+strconv.Itoa(int(resolveContext.RegisteredModules["gateway"].Module.Services(resolveContext.Context)[0].Port))),
			modules.Env("STARGATE_AUTH_CLIENT_ID", resolveContext.Stack.Spec.Stargate.AuthClientID),
			modules.Env("STARGATE_AUTH_CLIENT_SECRET", resolveContext.Stack.Spec.Stargate.AuthClientSecret),
			modules.Env("STARGATE_AUTH_ISSUER_URL", resolveContext.Stack.Spec.Stargate.AuthIssuerURL),
		)
	}
	modules.Register("stargate", modules.Module{
		Services: func(ctx modules.Context) modules.Services {
			return modules.Services{
				{
					AuthConfiguration: func(resolveContext modules.PrepareContext) stackv1beta3.ClientConfiguration {
						return stackv1beta3.NewClientConfiguration()
					},
					Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
						return modules.Container{
							Env:      stargateClientEnvVars(resolveContext),
							Image:    modules.GetImage("stargate", resolveContext.Versions.Spec.Stargate),
							Liveness: modules.LivenessDisable,
						}
					},
				},
			}
		},
	})
}

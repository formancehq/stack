package handlers

import (
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

func init() {
	modules.Register("wallets", modules.Module{
		Services: func(ctx modules.Context) modules.Services {
			return modules.Services{{
				HasVersionEndpoint: true,
				ExposeHTTP:         true,
				ListenEnvVar:       "LISTEN",
				AuthConfiguration: func(resolveContext modules.PrepareContext) stackv1beta3.ClientConfiguration {
					return stackv1beta3.NewClientConfiguration()
				},
				Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
					return modules.Container{
						Env: modules.ContainerEnv{
							modules.Env("STORAGE_POSTGRES_CONN_STRING", "$(POSTGRES_URI)"),
							modules.Env("STACK_CLIENT_ID", resolveContext.Stack.Status.StaticAuthClients["wallets"].ID),
							modules.Env("STACK_CLIENT_SECRET", resolveContext.Stack.Status.StaticAuthClients["wallets"].Secrets[0]),
						},
						Image:     modules.GetImage("wallets", resolveContext.Versions.Spec.Wallets),
						Resources: modules.ResourceSizeSmall(),
					}
				},
			}}
		},
	})
}

package handlers

import (
	"fmt"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
)

func init() {
	modules.Register("control", modules.Module{
		Versions: map[string]modules.Version{
			"v0.0.0": {
				Services: func(ctx modules.ModuleContext) modules.Services {
					return modules.Services{{
						Secured:     true,
						Port:        3000,
						ExposeHTTP:  true,
						Liveness:    modules.LivenessDisable,
						Annotations: ctx.Configuration.Spec.Services.Control.Annotations.Service,
						AuthConfiguration: func(resolveContext modules.ModuleContext) stackv1beta3.ClientConfiguration {
							return stackv1beta3.NewClientConfiguration().
								WithAdditionalScopes("profile", "email", "offline").
								WithRedirectUris(fmt.Sprintf("%s/auth/login", resolveContext.Stack.URL())).
								WithPostLogoutRedirectUris(fmt.Sprintf("%s/auth/destroy", resolveContext.Stack.URL()))
						},
						Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
							url := resolveContext.Stack.URL()
							if !resolveContext.HasVersionHigherOrEqual("v1.8.0") {
								url = fmt.Sprintf("%s/api", resolveContext.Stack.URL())
							}
							env := modules.ContainerEnv{
								modules.Env("API_URL", url),
								modules.Env("ENCRYPTION_KEY", "9h44y2ZqrDuUy5R9NGLA9hca7uRUr932"),
								modules.Env("ENCRYPTION_IV", "b6747T6eP9DnMvEw"),
								modules.Env("CLIENT_ID", resolveContext.Stack.Status.StaticAuthClients["control"].ID),
								modules.Env("CLIENT_SECRET", resolveContext.Stack.Status.StaticAuthClients["control"].Secrets[0]),
								modules.Env("REDIRECT_URI", resolveContext.Stack.URL()),
								modules.EnvFromBool("UNSECURE_COOKIES", resolveContext.Stack.Spec.Dev),
							}
							return modules.Container{
								Name:  "control",
								Image: modules.GetImage("control", resolveContext.Versions.Spec.Control),
								Env:   env,
								Resources: getResourcesWithDefault(
									resolveContext.Configuration.Spec.Services.Control.ResourceProperties,
									modules.ResourceSizeMedium(),
								),
							}
						},
					}}
				},
			},
		},
	})
}

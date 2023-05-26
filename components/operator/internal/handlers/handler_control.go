package handlers

import (
	"fmt"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	"github.com/rogpeppe/go-internal/semver"
)

func init() {
	modules.Register("control", modules.Module{
		Services: func(ctx modules.Context) modules.Services {
			return modules.Services{{
				Secured:    true,
				Port:       3000,
				ExposeHTTP: true,
				Liveness:   modules.LivenessDisable,
				AuthConfiguration: func(resolveContext modules.PrepareContext) stackv1beta3.ClientConfiguration {
					return stackv1beta3.NewClientConfiguration().
						WithAdditionalScopes("profile", "email", "offline").
						WithRedirectUris(fmt.Sprintf("%s/auth/login", resolveContext.Stack.URL())).
						WithPostLogoutRedirectUris(fmt.Sprintf("%s/auth/destroy", resolveContext.Stack.URL()))
				},
				Container: func(resolveContext modules.ContainerResolutionContext) modules.Container {
					tag := resolveContext.Versions.Spec.Control
					url := ""
					if semver.IsValid(tag) {
						version := semver.Compare(tag, "v1.8.0")
						if version == +1 {
							url = resolveContext.Stack.URL()
						}
						if version == 0 {
							url = resolveContext.Stack.URL()
						}
						if version == -1 {
							url = fmt.Sprintf("%s/api", resolveContext.Stack.URL())
						}
					}
					if url == "" {
						url = resolveContext.Stack.URL()
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
						Name:      "control",
						Image:     modules.GetImage("control", resolveContext.Versions.Spec.Control),
						Env:       env,
						Resources: modules.ResourceSizeMedium(),
					}
				},
			}}
		},
	})
}

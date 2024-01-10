package auths

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/stacks"
	v1 "k8s.io/api/core/v1"
	"strconv"
)

func EnvVars(ctx Context, stack *v1beta1.Stack, moduleName string, auth *v1beta1.AuthConfig) ([]v1.EnvVar, error) {
	ret := make([]v1.EnvVar, 0)

	hasAuth, err := stacks.HasDependency[*v1beta1.Auth](ctx, stack.Name)
	if err != nil {
		return nil, err
	}
	if !hasAuth {
		return ret, nil
	}

	url, err := URL(ctx, stack.Name)
	if err != nil {
		return nil, err
	}

	ret = append(ret,
		Env("AUTH_ENABLED", "true"),
		Env("AUTH_ISSUER", fmt.Sprintf("%s/api/auth", url)),
	)

	if auth != nil {
		if auth.ReadKeySetMaxRetries != 0 {
			ret = append(ret,
				Env("AUTH_READ_KEY_SET_MAX_RETRIES", strconv.Itoa(auth.ReadKeySetMaxRetries)),
			)
		}

		if auth.CheckScopes {
			ret = append(ret,
				Env("AUTH_CHECK_SCOPES", "true"),
				Env("AUTH_SERVICE", moduleName),
			)
		}
	}

	return ret, nil
}

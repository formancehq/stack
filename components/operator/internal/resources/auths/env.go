package auths

import (
	"fmt"
	"strconv"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"
	v1 "k8s.io/api/core/v1"
)

func ProtectedEnvVars(ctx Context, stack *v1beta1.Stack, moduleName string, auth *v1beta1.AuthConfig) ([]v1.EnvVar, error) {
	return ProtectedAPIEnvVarsWithPrefix(ctx, stack, moduleName, auth, "")
}

func ProtectedAPIEnvVarsWithPrefix(ctx Context, stack *v1beta1.Stack, moduleName string, auth *v1beta1.AuthConfig, prefix string) ([]v1.EnvVar, error) {
	ret := make([]v1.EnvVar, 0)

	hasAuth, err := HasDependency(ctx, stack.Name, &v1beta1.Auth{})
	if err != nil {
		return nil, err
	}
	if !hasAuth {
		return ret, nil
	}

	url, err := getUrl(ctx, stack.Name)
	if err != nil {
		return nil, err
	}

	ret = append(ret,
		Env(fmt.Sprintf("%sAUTH_ENABLED", prefix), "true"),
		Env(fmt.Sprintf("%sAUTH_ISSUER", prefix), fmt.Sprintf(fmt.Sprintf("%s/api/auth", url))),
	)

	if auth != nil {
		if auth.ReadKeySetMaxRetries != 0 {
			ret = append(ret,
				Env(fmt.Sprintf("%sAUTH_READ_KEY_SET_MAX_RETRIES", prefix), strconv.Itoa(auth.ReadKeySetMaxRetries)),
			)
		}

		if auth.CheckScopes {
			ret = append(ret,
				Env(fmt.Sprintf("%sAUTH_CHECK_SCOPES", prefix), "true"),
				Env(fmt.Sprintf("%sAUTH_SERVICE", prefix), moduleName),
			)
		}
	}

	return ret, nil
}

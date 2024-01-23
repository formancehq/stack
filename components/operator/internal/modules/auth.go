package modules

import (
	"strconv"

	"github.com/formancehq/operator/apis/stack/v1beta3"
)

func AuthEnvVars(stackURL, moduleName string, auth *v1beta3.AuthConfig) ContainerEnv {
	ret := ContainerEnv{}

	ret = ret.Append(
		Env("AUTH_ENABLED", "false"),
	)

	if auth != nil {
		if auth.ReadKeySetMaxRetries != 0 {
			ret = ret.Append(
				Env("AUTH_READ_KEY_SET_MAX_RETRIES", strconv.Itoa(auth.ReadKeySetMaxRetries)),
			)
		}

		if auth.CheckScopes {
			ret = ret.Append(
				Env("AUTH_CHECK_SCOPES", "true"),
				Env("AUTH_SERVICE", moduleName),
			)
		}
	}

	return ret
}

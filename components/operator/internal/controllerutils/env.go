package controllerutils

import (
	"fmt"

	"github.com/formancehq/operator/internal/collectionutils"
)

func EnvVarPlaceholder(key, prefix string) string {
	return fmt.Sprintf("$(%s%s)", prefix, key)
}

func ComputeEnvVar(prefix, format string, keys ...string) string {
	return fmt.Sprintf(format,
		collectionutils.Map(keys, func(key string) any {
			return EnvVarPlaceholder(key, prefix)
		})...,
	)
}

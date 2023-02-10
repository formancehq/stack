package controllerutils

import (
	"fmt"

	"github.com/formancehq/operator/internal/collectionutils"
	v1 "k8s.io/api/core/v1"
)

func EnvWithPrefix(prefix, key, value string) v1.EnvVar {
	return v1.EnvVar{
		Name:  prefix + key,
		Value: value,
	}
}

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

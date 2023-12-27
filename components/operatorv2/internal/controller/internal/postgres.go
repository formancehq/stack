package internal

import (
	"fmt"

	"github.com/formancehq/operator/v2/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func PostgresEnvVars(c v1beta1.DatabaseConfigurationSpec, dbName string) []corev1.EnvVar {
	ret := []corev1.EnvVar{
		Env("POSTGRES_HOST", c.Host),
		Env("POSTGRES_PORT", fmt.Sprint(c.Port)),
		Env("POSTGRES_DATABASE", dbName),
	}
	if c.Username != "" || c.CredentialsFromSecret != "" {
		if c.Username != "" {
			ret = append(ret,
				Env("POSTGRES_USERNAME", c.Username),
				Env("POSTGRES_PASSWORD", c.Password),
			)
		} else {
			ret = append(ret,
				EnvFromSecret("POSTGRES_USERNAME", c.CredentialsFromSecret, "username"),
				EnvFromSecret("POSTGRES_PASSWORD", c.CredentialsFromSecret, "password"),
			)
		}
		ret = append(ret,
			Env("POSTGRES_NO_DATABASE_URI", ComputeEnvVar("postgresql://%s:%s@%s:%s",
				"POSTGRES_USERNAME",
				"POSTGRES_PASSWORD",
				"POSTGRES_HOST",
				"POSTGRES_PORT",
			)),
		)
	} else {
		ret = append(ret,
			Env("POSTGRES_NO_DATABASE_URI", ComputeEnvVar("postgresql://%s:%s",
				"POSTGRES_HOST",
				"POSTGRES_PORT",
			)),
		)
	}

	fmt := "%s/%s"
	if c.DisableSSLMode {
		fmt += "?sslmode=disable"
	}
	ret = append(ret,
		Env("POSTGRES_URI", ComputeEnvVar(fmt, "POSTGRES_NO_DATABASE_URI", "POSTGRES_DATABASE")),
	)

	return ret
}

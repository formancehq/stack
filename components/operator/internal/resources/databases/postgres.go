package databases

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	corev1 "k8s.io/api/core/v1"
)

func PostgresEnvVars(c v1beta1.DatabaseConfigurationSpec, dbName string) []corev1.EnvVar {
	ret := []corev1.EnvVar{
		core.Env("POSTGRES_HOST", c.Host),
		core.Env("POSTGRES_PORT", fmt.Sprint(c.Port)),
		core.Env("POSTGRES_DATABASE", dbName),
	}
	if c.Username != "" || c.CredentialsFromSecret != "" {
		if c.Username != "" {
			ret = append(ret,
				core.Env("POSTGRES_USERNAME", c.Username),
				core.Env("POSTGRES_PASSWORD", c.Password),
			)
		} else {
			ret = append(ret,
				core.EnvFromSecret("POSTGRES_USERNAME", c.CredentialsFromSecret, "username"),
				core.EnvFromSecret("POSTGRES_PASSWORD", c.CredentialsFromSecret, "password"),
			)
		}
		ret = append(ret,
			core.Env("POSTGRES_NO_DATABASE_URI", core.ComputeEnvVar("postgresql://%s:%s@%s:%s",
				"POSTGRES_USERNAME",
				"POSTGRES_PASSWORD",
				"POSTGRES_HOST",
				"POSTGRES_PORT",
			)),
		)
	} else {
		ret = append(ret,
			core.Env("POSTGRES_NO_DATABASE_URI", core.ComputeEnvVar("postgresql://%s:%s",
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
		core.Env("POSTGRES_URI", core.ComputeEnvVar(fmt, "POSTGRES_NO_DATABASE_URI", "POSTGRES_DATABASE")),
	)

	return ret
}

package databases

import (
	"fmt"
	"github.com/formancehq/operator/v2/internal/utils"

	"github.com/formancehq/operator/v2/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func PostgresEnvVars(c v1beta1.DatabaseConfigurationSpec, dbName string) []corev1.EnvVar {
	ret := []corev1.EnvVar{
		utils.Env("POSTGRES_HOST", c.Host),
		utils.Env("POSTGRES_PORT", fmt.Sprint(c.Port)),
		utils.Env("POSTGRES_DATABASE", dbName),
	}
	if c.Username != "" || c.CredentialsFromSecret != "" {
		if c.Username != "" {
			ret = append(ret,
				utils.Env("POSTGRES_USERNAME", c.Username),
				utils.Env("POSTGRES_PASSWORD", c.Password),
			)
		} else {
			ret = append(ret,
				utils.EnvFromSecret("POSTGRES_USERNAME", c.CredentialsFromSecret, "username"),
				utils.EnvFromSecret("POSTGRES_PASSWORD", c.CredentialsFromSecret, "password"),
			)
		}
		ret = append(ret,
			utils.Env("POSTGRES_NO_DATABASE_URI", utils.ComputeEnvVar("postgresql://%s:%s@%s:%s",
				"POSTGRES_USERNAME",
				"POSTGRES_PASSWORD",
				"POSTGRES_HOST",
				"POSTGRES_PORT",
			)),
		)
	} else {
		ret = append(ret,
			utils.Env("POSTGRES_NO_DATABASE_URI", utils.ComputeEnvVar("postgresql://%s:%s",
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
		utils.Env("POSTGRES_URI", utils.ComputeEnvVar(fmt, "POSTGRES_NO_DATABASE_URI", "POSTGRES_DATABASE")),
	)

	return ret
}

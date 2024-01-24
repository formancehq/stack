package databases

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"

	"github.com/formancehq/operator/internal/core"
	corev1 "k8s.io/api/core/v1"
)

func PostgresEnvVars(c v1beta1.DatabaseConfiguration, dbName string) []corev1.EnvVar {
	return PostgresEnvVarsWithPrefix(c, dbName, "")
}

func PostgresEnvVarsWithPrefix(c v1beta1.DatabaseConfiguration, dbName, prefix string) []corev1.EnvVar {
	ret := []corev1.EnvVar{
		core.Env(fmt.Sprintf("%sPOSTGRES_HOST", prefix), c.Host),
		core.Env(fmt.Sprintf("%sPOSTGRES_PORT", prefix), fmt.Sprint(c.Port)),
		core.Env(fmt.Sprintf("%sPOSTGRES_DATABASE", prefix), dbName),
	}
	if c.Username != "" || c.CredentialsFromSecret != "" {
		if c.Username != "" {
			ret = append(ret,
				core.Env(fmt.Sprintf("%sPOSTGRES_USERNAME", prefix), c.Username),
				core.Env(fmt.Sprintf("%sPOSTGRES_PASSWORD", prefix), c.Password),
			)
		} else {
			ret = append(ret,
				core.EnvFromSecret(fmt.Sprintf("%sPOSTGRES_USERNAME", prefix), c.CredentialsFromSecret, "username"),
				core.EnvFromSecret(fmt.Sprintf("%sPOSTGRES_PASSWORD", prefix), c.CredentialsFromSecret, "password"),
			)
		}
		ret = append(ret,
			core.Env(fmt.Sprintf("%sPOSTGRES_NO_DATABASE_URI", prefix), core.ComputeEnvVar("postgresql://%s:%s@%s:%s",
				fmt.Sprintf("%sPOSTGRES_USERNAME", prefix),
				fmt.Sprintf("%sPOSTGRES_PASSWORD", prefix),
				fmt.Sprintf("%sPOSTGRES_HOST", prefix),
				fmt.Sprintf("%sPOSTGRES_PORT", prefix),
			)),
		)
	} else {
		ret = append(ret,
			core.Env(fmt.Sprintf("%sPOSTGRES_NO_DATABASE_URI", prefix), core.ComputeEnvVar("postgresql://%s:%s",
				fmt.Sprintf("%sPOSTGRES_HOST", prefix),
				fmt.Sprintf("%sPOSTGRES_PORT", prefix),
			)),
		)
	}

	f := "%s/%s"
	if c.DisableSSLMode {
		f += "?sslmode=disable"
	}
	ret = append(ret,
		core.Env(fmt.Sprintf("%sPOSTGRES_URI", prefix), core.ComputeEnvVar(f,
			fmt.Sprintf("%sPOSTGRES_NO_DATABASE_URI", prefix),
			fmt.Sprintf("%sPOSTGRES_DATABASE", prefix))),
	)

	return ret
}

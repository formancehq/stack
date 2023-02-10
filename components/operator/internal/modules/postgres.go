package modules

import (
	"fmt"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllerutils"
)

func DefaultPostgresEnvVarsWithPrefix(c v1beta3.PostgresConfig, dbName, prefix string) ContainerEnv {

	withPrefix := func(v string) string {
		return prefix + v
	}

	ret := ContainerEnv{
		Env(withPrefix("POSTGRES_HOST"), c.Host),
		Env(withPrefix("POSTGRES_PORT"), fmt.Sprint(c.Port)),
		Env(withPrefix("POSTGRES_DATABASE"), dbName),
	}
	if c.Username != "" {
		ret = ret.Append(
			Env(withPrefix("POSTGRES_USERNAME"), c.Username),
			Env(withPrefix("POSTGRES_PASSWORD"), c.Password),
			Env(withPrefix("POSTGRES_NO_DATABASE_URI"), controllerutils.ComputeEnvVar(prefix, "postgresql://%s:%s@%s:%s",
				"POSTGRES_USERNAME",
				"POSTGRES_PASSWORD",
				"POSTGRES_HOST",
				"POSTGRES_PORT",
			)),
		)
	} else {
		ret = ret.Append(
			Env(withPrefix("POSTGRES_NO_DATABASE_URI"), controllerutils.ComputeEnvVar(prefix, "postgresql://%s:%s",
				"POSTGRES_HOST",
				"POSTGRES_PORT",
			)),
		)
	}

	fmt := "%s/%s"
	if c.DisableSSLMode {
		fmt += "?sslmode=disable"
	}
	ret = ret.Append(
		Env(prefix+"POSTGRES_URI", controllerutils.ComputeEnvVar(prefix, fmt, "POSTGRES_NO_DATABASE_URI", "POSTGRES_DATABASE")),
	)

	return ret
}

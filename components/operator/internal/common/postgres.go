package common

import (
	"fmt"
	"strings"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/collectionutils"
	"github.com/formancehq/operator/internal/controllerutils"
	corev1 "k8s.io/api/core/v1"
)

func DefaultPostgresEnvVarsLegacy(c v1beta3.PostgresConfig, dbName string, prefix string) []corev1.EnvVar {
	return defaultPostgresEnvVarsWithDiscriminatorLegacy(c, dbName, prefix, "")
}

func defaultPostgresEnvVarsWithDiscriminatorLegacy(c v1beta3.PostgresConfig, dbName string, prefix, discriminator string) []corev1.EnvVar {

	discriminator = strings.ToUpper(discriminator)
	withDiscriminator := func(v string) string {
		if discriminator == "" {
			return v
		}
		return fmt.Sprintf("%s_%s", v, discriminator)
	}

	ret := make([]corev1.EnvVar, 0)
	ret = append(ret, controllerutils.EnvWithPrefix(prefix, withDiscriminator("POSTGRES_HOST"), c.Host))
	ret = append(ret, controllerutils.EnvWithPrefix(prefix, withDiscriminator("POSTGRES_PORT"), fmt.Sprint(c.Port)))
	ret = append(ret, controllerutils.EnvWithPrefix(prefix, withDiscriminator("POSTGRES_DATABASE"), dbName))
	if c.Username != "" {
		ret = append(ret, controllerutils.EnvWithPrefix(prefix, withDiscriminator("POSTGRES_USERNAME"), c.Username))
		ret = append(ret, controllerutils.EnvWithPrefix(prefix, withDiscriminator("POSTGRES_PASSWORD"), c.Password))

		ret = append(ret, controllerutils.EnvWithPrefix(prefix, withDiscriminator("POSTGRES_NO_DATABASE_URI"),
			controllerutils.ComputeEnvVar(prefix, "postgresql://%s:%s@%s:%s",
				withDiscriminator("POSTGRES_USERNAME"),
				withDiscriminator("POSTGRES_PASSWORD"),
				withDiscriminator("POSTGRES_HOST"),
				withDiscriminator("POSTGRES_PORT"),
			),
		))
	} else {
		ret = append(ret, controllerutils.EnvWithPrefix(prefix, withDiscriminator("POSTGRES_NO_DATABASE_URI"),
			controllerutils.ComputeEnvVar(prefix, "postgresql://%s:%s", withDiscriminator("POSTGRES_HOST"), withDiscriminator("POSTGRES_PORT")),
		))
	}
	fmt := "%s/%s"
	if c.DisableSSLMode {
		fmt += "?sslmode=disable"
	}
	ret = append(ret,
		controllerutils.EnvWithPrefix(prefix, withDiscriminator("POSTGRES_URI"),
			controllerutils.ComputeEnvVar(prefix, fmt, withDiscriminator("POSTGRES_NO_DATABASE_URI"), withDiscriminator("POSTGRES_DATABASE")),
		),
	)

	return ret
}

func CreateDatabaseInitContainer(c v1beta3.PostgresConfig, dbName, prefix string) corev1.Container {
	return corev1.Container{
		Name:            "init-db-" + dbName,
		Image:           "postgres:13",
		ImagePullPolicy: corev1.PullIfNotPresent,
		Env:             collectionutils.NewArray[corev1.EnvVar]().Append(DefaultPostgresEnvVarsLegacy(c, dbName, prefix)...),
		Command: []string{
			"sh",
			"-c",
			fmt.Sprintf(`psql -Atx ${%sPOSTGRES_NO_DATABASE_URI}/postgres -c "SELECT 1 FROM pg_database WHERE datname = '${%sPOSTGRES_DATABASE}'" | grep -q 1 && echo "Base already exists" || psql -Atx ${%sPOSTGRES_NO_DATABASE_URI}/postgres -c "CREATE DATABASE \"${%sPOSTGRES_DATABASE}\""`,
				prefix, prefix, prefix, prefix,
			),
		},
	}
}

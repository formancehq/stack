package databases

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/resources/settings"
	"net/url"

	"github.com/formancehq/operator/internal/core"
	corev1 "k8s.io/api/core/v1"
)

func GetPostgresEnvVars(db *v1beta1.Database) []corev1.EnvVar {
	return PostgresEnvVarsWithPrefix(db, "")
}

func PostgresEnvVarsWithPrefix(database *v1beta1.Database, prefix string) []corev1.EnvVar {
	databaseURL := database.Status.URL()
	ret := []corev1.EnvVar{
		core.Env(fmt.Sprintf("%sPOSTGRES_HOST", prefix), databaseURL.Hostname()),
		core.Env(fmt.Sprintf("%sPOSTGRES_PORT", prefix), databaseURL.Port()),
		core.Env(fmt.Sprintf("%sPOSTGRES_DATABASE", prefix), database.Status.Database),
	}
	if databaseURL.User.Username() != "" || databaseURL.Query().Get("secret") != "" {
		if databaseURL.User.Username() != "" {
			password, _ := databaseURL.User.Password()
			ret = append(ret,
				core.Env(fmt.Sprintf("%sPOSTGRES_USERNAME", prefix), databaseURL.User.Username()),
				core.Env(fmt.Sprintf("%sPOSTGRES_PASSWORD", prefix), password),
			)
		} else {
			secret := fmt.Sprintf("%s-%s", database.Name, databaseURL.Query().Get("secret"))
			ret = append(ret,
				core.EnvFromSecret(fmt.Sprintf("%sPOSTGRES_USERNAME", prefix), secret, "username"),
				core.EnvFromSecret(fmt.Sprintf("%sPOSTGRES_PASSWORD", prefix), secret, "password"),
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
	if isDisabledSSLMode(databaseURL) {
		f += "?sslmode=disable"
	}
	ret = append(ret,
		core.Env(fmt.Sprintf("%sPOSTGRES_URI", prefix), core.ComputeEnvVar(f,
			fmt.Sprintf("%sPOSTGRES_NO_DATABASE_URI", prefix),
			fmt.Sprintf("%sPOSTGRES_DATABASE", prefix))),
	)

	return ret
}

func isDisabledSSLMode(url *url.URL) bool {
	return settings.IsTrue(url.Query().Get("disableSSLMode"))
}

package databases

import (
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	corev1 "k8s.io/api/core/v1"
)

func GetPostgresEnvVars(ctx core.Context, stack *v1beta1.Stack, db *v1beta1.Database) ([]corev1.EnvVar, error) {
	return PostgresEnvVarsWithPrefix(ctx, stack, db, "")
}

func PostgresEnvVarsWithPrefix(ctx core.Context, stack *v1beta1.Stack, database *v1beta1.Database, prefix string) ([]corev1.EnvVar, error) {
	ret := []corev1.EnvVar{
		core.Env(fmt.Sprintf("%sPOSTGRES_HOST", prefix), database.Status.URI.Hostname()),
		core.Env(fmt.Sprintf("%sPOSTGRES_PORT", prefix), database.Status.URI.Port()),
		core.Env(fmt.Sprintf("%sPOSTGRES_DATABASE", prefix), database.Status.Database),
	}
	if database.Status.URI.User.Username() != "" || database.Status.URI.Query().Get("secret") != "" {
		if database.Status.URI.User.Username() != "" {
			password, _ := database.Status.URI.User.Password()
			ret = append(ret,
				core.Env(fmt.Sprintf("%sPOSTGRES_USERNAME", prefix), database.Status.URI.User.Username()),
				core.Env(fmt.Sprintf("%sPOSTGRES_PASSWORD", prefix), password),
			)
		} else {
			secret := database.Status.URI.Query().Get("secret")
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

	awsRole, err := settings.GetAWSServiceAccount(ctx, stack.Name)
	if err != nil {
		return nil, err
	}

	if awsRole != "" {
		ret = append(ret, core.Env(fmt.Sprintf("%sPOSTGRES_AWS_ENABLE_IAM", prefix), "true"))
	}

	f := "%s/%s"
	if settings.IsTrue(database.Status.URI.Query().Get("disableSSLMode")) {
		f += "?sslmode=disable"
	}
	ret = append(ret,
		core.Env(fmt.Sprintf("%sPOSTGRES_URI", prefix), core.ComputeEnvVar(f,
			fmt.Sprintf("%sPOSTGRES_NO_DATABASE_URI", prefix),
			fmt.Sprintf("%sPOSTGRES_DATABASE", prefix))),
	)

	return ret, nil
}

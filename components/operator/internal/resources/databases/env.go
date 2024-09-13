package databases

import (
	"fmt"
	"strconv"
	"time"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/pkg/errors"
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

	config, err := settings.GetAs[connectionPoolConfiguration](ctx, stack.Name, "modules", database.Spec.Service, "database", "connection-pool")
	if err != nil {
		return nil, err
	}

	if config.MaxIdle != "" {
		_, err := strconv.ParseUint(config.MaxIdle, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse max idle value")
		}
		ret = append(ret, core.Env(fmt.Sprintf("%sPOSTGRES_MAX_IDLE_CONNS", prefix), config.MaxIdle))
	}
	if config.MaxIdleTime != "" {
		_, err := time.ParseDuration(config.MaxIdleTime)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse max idle time value")
		}
		ret = append(ret, core.Env(fmt.Sprintf("%sPOSTGRES_CONN_MAX_IDLE_TIME", prefix), config.MaxIdleTime))
	}
	if config.MaxOpen != "" {
		_, err := strconv.ParseUint(config.MaxOpen, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse max open conns value")
		}
		ret = append(ret, core.Env(fmt.Sprintf("%sPOSTGRES_CONN_MAX_OPEN_CONNS", prefix), config.MaxOpen))
	}

	return ret, nil
}

type connectionPoolConfiguration struct {
	MaxIdle     string `json:"max-idle,omitempty"`
	MaxIdleTime string `json:"max-idle-time,omitempty"`
	MaxOpen     string `json:"max-open,omitempty"`
}

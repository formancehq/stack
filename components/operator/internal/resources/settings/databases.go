package settings

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
)

func NewHostSetting(name, value string, stacks ...string) *v1beta1.Settings {
	return New(name, "databases.host", value, stacks...)
}

func FindDatabaseConfiguration(ctx core.Context, database *v1beta1.Database) (*v1beta1.DatabaseConfiguration, error) {
	host, err := GetString(ctx, database.Spec.Stack, "databases", database.Spec.Service, "host")
	if err != nil {
		return nil, err
	}
	if host == nil {
		host, err = GetString(ctx, database.Spec.Stack, "databases", "host")
		if err != nil {
			return nil, err
		}
	}

	if host == nil {
		return nil, errors.New("missing database host")
	}

	port, err := GetInt64(ctx, database.Spec.Stack, "databases", database.Spec.Service, "port")
	if err != nil {
		return nil, err
	}
	if port == nil {
		port, err = GetInt64(ctx, database.Spec.Stack, "databases", "port")
		if err != nil {
			return nil, err
		}
	}
	if port == nil {
		port = pointer.For(int64(5432)) // default postgres port
	}

	username, err := GetString(ctx, database.Spec.Stack, "databases", database.Spec.Service, "username")
	if err != nil {
		return nil, err
	}
	if username == nil {
		username, err = GetString(ctx, database.Spec.Stack, "databases", "username")
		if err != nil {
			return nil, err
		}
	}

	password, err := GetString(ctx, database.Spec.Stack, "databases", database.Spec.Service, "password")
	if err != nil {
		return nil, err
	}
	if password == nil {
		password, err = GetString(ctx, database.Spec.Stack, "databases", "password")
		if err != nil {
			return nil, err
		}
	}

	credentialsFromSecret, err := GetString(ctx, database.Spec.Stack, "databases", database.Spec.Service, "credentials-from-secret")
	if err != nil {
		return nil, err
	}
	if credentialsFromSecret == nil {
		credentialsFromSecret, err = GetString(ctx, database.Spec.Stack, "databases", "secret")
		if err != nil {
			return nil, err
		}
	}

	disableSSLMode, err := GetBool(ctx, database.Spec.Stack, "databases", database.Spec.Service, "ssl", "disable")
	if err != nil {
		return nil, err
	}
	if disableSSLMode == nil {
		disableSSLMode, err = GetBool(ctx, database.Spec.Stack, "databases", "disable-ssl-mode")
		if err != nil {
			return nil, err
		}
	}

	return &v1beta1.DatabaseConfiguration{
		Port:                  int(*port),
		Host:                  *host,
		Username:              ValueOrDefault(username, ""),
		Password:              ValueOrDefault(password, ""),
		CredentialsFromSecret: ValueOrDefault(credentialsFromSecret, ""),
		DisableSSLMode:        ValueOrDefault(disableSSLMode, false),
	}, nil
}

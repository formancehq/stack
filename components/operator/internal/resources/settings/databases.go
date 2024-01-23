package settings

import (
	"fmt"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"net/url"
	"strconv"
	"strings"
)

func NewHostSetting(...any) *v1beta1.Settings {
	return nil
}

func FindDatabaseConfiguration(ctx core.Context, database *v1beta1.Database) (*v1beta1.DatabaseConfiguration, error) {

	uri, err := GetString(ctx, database.Spec.Stack, "postgres", database.Spec.Service, "uri")
	if err != nil {
		return nil, err
	}
	if uri == nil {
		return nil, nil
	}

	parsedUrl, err := url.Parse(*uri)
	if err != nil {
		return nil, err
	}

	if parsedUrl.Scheme != "postgresql" {
		return nil, fmt.Errorf("invalid postgres uri: %s", *uri)
	}

	port := uint64(5432)
	if parsedUrl.Port() != "" {
		port, err = strconv.ParseUint(parsedUrl.Port(), 10, 16)
		if err != nil {
			return nil, err
		}
	}

	password, _ := parsedUrl.User.Password()

	return &v1beta1.DatabaseConfiguration{
		Port:                  int(port),
		Host:                  parsedUrl.Hostname(),
		Username:              parsedUrl.User.Username(),
		Password:              password,
		CredentialsFromSecret: parsedUrl.Query().Get("credentialsFromSecret"),
		DisableSSLMode:        strings.ToLower(parsedUrl.Query().Get("disableSSLMode")) == "true" || parsedUrl.Query().Get("disableSSLMode") == "1",
	}, nil
}

func isTrue(v string) bool {
	return strings.ToLower(v) == "true" || v == "1"
}

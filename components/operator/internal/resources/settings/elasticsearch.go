package settings

import (
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
)

func FindElasticSearchConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.ElasticSearchConfiguration, error) {
	elasticSearchHost, err := RequireString(ctx, stack.Name, "elasticsearch.host")
	if err != nil {
		return nil, err
	}

	elasticSearchScheme, err := GetStringOrDefault(ctx, stack.Name, "https", "elasticsearch.scheme")
	if err != nil {
		return nil, err
	}

	elasticSearchPort, err := GetUInt16OrDefault(ctx, stack.Name, 9200, "elasticsearch.port")
	if err != nil {
		return nil, err
	}

	elasticSearchTLSEnabled, err := GetBoolOrFalse(ctx, stack.Name, "elasticsearch.tls.enabled")
	if err != nil {
		return nil, err
	}

	elasticSearchTLSSkipCertVerify, err := GetBoolOrFalse(ctx, stack.Name, "elasticsearch.tls.skip-cert-verify")
	if err != nil {
		return nil, err
	}

	var basicAuth *v1beta1.ElasticSearchBasicAuthConfig
	basicAuthEnabled, err := GetBoolOrFalse(ctx, stack.Name, "elasticsearch.basic-auth.enabled")
	if basicAuthEnabled {
		elasticSearchBasicAuthUsername, err := GetStringOrEmpty(ctx, stack.Name, "elasticsearch.basic-auth.username")
		if err != nil {
			return nil, err
		}

		elasticSearchBasicAuthPassword, err := GetStringOrEmpty(ctx, stack.Name, "elasticsearch.basic-auth.password")
		if err != nil {
			return nil, err
		}

		elasticSearchBasicAuthSecret, err := GetStringOrEmpty(ctx, stack.Name, "elasticsearch.basic-auth.secret")
		if err != nil {
			return nil, err
		}

		basicAuth = &v1beta1.ElasticSearchBasicAuthConfig{
			Username:   elasticSearchBasicAuthUsername,
			Password:   elasticSearchBasicAuthPassword,
			SecretName: elasticSearchBasicAuthSecret,
		}
	}

	return &v1beta1.ElasticSearchConfiguration{
		Scheme: elasticSearchScheme,
		Host:   elasticSearchHost,
		Port:   elasticSearchPort,
		TLS: v1beta1.ElasticSearchTLSConfig{
			Enabled:        elasticSearchTLSEnabled,
			SkipCertVerify: elasticSearchTLSSkipCertVerify,
		},
		BasicAuth: basicAuth,
	}, nil
}

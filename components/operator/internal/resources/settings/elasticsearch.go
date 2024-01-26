package settings

import (
	"net/url"
	"strconv"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
)

func FindElasticSearchConfiguration(ctx core.Context, stack *v1beta1.Stack) (*v1beta1.ElasticSearchConfiguration, error) {
	elasticSearchDSN, err := GetString(ctx, stack.Name, "elasticsearch.dsn")
	if err != nil {
		return nil, err
	}
	if elasticSearchDSN == nil {
		return nil, nil
	}

	elasticSearchURI, err := url.Parse(*elasticSearchDSN)
	if err != nil {
		return nil, err
	}

	port := uint64(9200)
	if elasticSearchURI.Port() != "" {
		port, err = strconv.ParseUint(elasticSearchURI.Port(), 10, 16)
		if err != nil {
			return nil, err
		}
	}

	var basicAuth *v1beta1.ElasticSearchBasicAuthConfig
	if elasticSearchURI.User.Username() != "" {
		password, _ := elasticSearchURI.User.Password()

		basicAuth = &v1beta1.ElasticSearchBasicAuthConfig{
			Username: elasticSearchURI.User.Username(),
			Password: password,
		}
	} else {
		basicAuth = &v1beta1.ElasticSearchBasicAuthConfig{
			SecretName: elasticSearchURI.Query().Get("secret"),
		}
	}

	return &v1beta1.ElasticSearchConfiguration{
		Scheme: elasticSearchURI.Scheme,
		Host:   elasticSearchURI.Hostname(),
		Port:   uint16(port),
		TLS: v1beta1.ElasticSearchTLSConfig{
			Enabled:        IsTrue(elasticSearchURI.Query().Get("tls")),
			SkipCertVerify: IsTrue(elasticSearchURI.Query().Get("skipCertVerify")),
		},
		BasicAuth: basicAuth,
	}, nil
}

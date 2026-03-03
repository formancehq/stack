package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zitadel/oidc/v2/pkg/client"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/op"
	"go.uber.org/fx"
)

type ModuleConfig struct {
	Enabled              bool
	Issuers              []string
	ReadKeySetMaxRetries int
	CheckScopes          bool
	Service              string

	// Deprecated: use Issuers instead.
	Issuer string
}

func (cfg ModuleConfig) resolveIssuers() []string {
	issuers := cfg.Issuers
	if cfg.Issuer != "" {
		found := false
		for _, iss := range issuers {
			if iss == cfg.Issuer {
				found = true
				break
			}
		}
		if !found {
			issuers = append(issuers, cfg.Issuer)
		}
	}
	return issuers
}

func Module(cfg ModuleConfig) fx.Option {
	options := make([]fx.Option, 0)

	issuers := cfg.resolveIssuers()

	if cfg.Enabled && len(issuers) == 0 {
		return fx.Error(errors.New("auth is enabled but no issuers are configured"))
	}

	options = append(options,
		fx.Provide(func() Auth {
			return NewNoAuth()
		}),
	)

	if cfg.Enabled {
		options = append(options,
			fx.Decorate(func(logger logging.Logger) (Auth, error) {
				retryClient := retryablehttp.NewClient()
				retryClient.RetryMax = cfg.ReadKeySetMaxRetries
				discoveryHTTPClient := retryClient.StandardClient()

				verifiers := make(map[string]op.AccessTokenVerifier, len(issuers))
				for _, issuer := range issuers {
					discovery, err := client.Discover(issuer, discoveryHTTPClient)
					if err != nil {
						return nil, err
					}
					keySet := rp.NewRemoteKeySet(&http.Client{Timeout: 10 * time.Second}, discovery.JwksURI)
					verifiers[issuer] = op.NewAccessTokenVerifier(issuer, keySet)
				}
				return newJWTAuth(
					logger,
					cfg.ReadKeySetMaxRetries,
					verifiers,
					cfg.Service,
					cfg.CheckScopes,
				), nil
			}),
		)
	}

	return fx.Options(options...)
}

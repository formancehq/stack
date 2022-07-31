package delegatedauth

import (
	"github.com/coreos/go-oidc"
	"go.uber.org/fx"
)

func Module(issuer, clientID, clientSecret, redirectURL string) fx.Option {
	return fx.Options(
		fx.Provide(ProvideDelegatedOIDCProvider),
		fx.Supply(Issuer(issuer)),
		fx.Provide(func(provider *oidc.Provider) OAuth2Config {
			return OAuth2Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				RedirectURL:  redirectURL,
				Endpoint:     provider.Endpoint(),
				Scopes:       []string{oidc.ScopeOpenID, "email"},
			}
		}),
	)
}

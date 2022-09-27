package delegatedauth

import (
	"github.com/zitadel/oidc/pkg/client/rp"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg Config) (rp.RelyingParty, error) {
			return rp.NewRelyingPartyOIDC(cfg.Issuer, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL, []string{"openid email"})
		}),
	)
}

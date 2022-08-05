package delegatedauth

import (
	"github.com/zitadel/oidc/pkg/client/rp"
	"go.uber.org/fx"
)

type Config struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func Module(cfg Config) fx.Option {
	return fx.Options(
		fx.Provide(func() (rp.RelyingParty, error) {
			return rp.NewRelyingPartyOIDC(cfg.Issuer, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL, []string{"openid email"})
		}),
	)
}

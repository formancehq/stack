package delegatedauth

import (
	"context"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg Config, httpClient *http.Client) (rp.RelyingParty, error) {
			return rp.NewRelyingPartyOIDC(context.Background(), cfg.Issuer, cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL, []string{"openid email"},
				rp.WithHTTPClient(httpClient))
		}),
	)
}

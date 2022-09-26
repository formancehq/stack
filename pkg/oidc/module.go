package oidc

import (
	"context"
	"crypto/rsa"
	"net/url"

	auth "github.com/formancehq/auth/pkg"
	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
	"go.uber.org/fx"
)

func Module(privateKey *rsa.PrivateKey, baseUrl *url.URL) fx.Option {
	return fx.Options(
		fx.Invoke(func(router *mux.Router, provider op.OpenIDProvider, storage Storage, relyingParty rp.RelyingParty) {
			AddRoutes(router, provider, storage, relyingParty, baseUrl)
		}),
		fx.Provide(fx.Annotate(func(storage Storage, relyingParty rp.RelyingParty, opts []auth.ClientOptions) *storageFacade {
			var staticClients []auth.Client
			for _, c := range opts {
				staticClients = append(staticClients, *auth.NewClient(c))
			}
			return NewStorageFacade(storage, relyingParty, privateKey, staticClients...)
		}, fx.As(new(op.Storage)))),
		fx.Provide(func(storage op.Storage) (op.OpenIDProvider, error) {
			return NewOpenIDProvider(context.TODO(), storage, baseUrl.String())
		}),
	)
}

package oidc

import (
	"context"
	"crypto/rsa"
	"net/url"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/delegatedauth"
	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
	"go.uber.org/fx"
)

func Module(privateKey *rsa.PrivateKey, baseUrl *url.URL, staticClients ...auth.StaticClient) fx.Option {
	return fx.Options(
		fx.Invoke(func(router *mux.Router, provider op.OpenIDProvider, storage Storage, relyingParty rp.RelyingParty) {
			AddRoutes(router, provider, storage, relyingParty, baseUrl)
		}),
		fx.Provide(fx.Annotate(func(storage Storage, relyingParty rp.RelyingParty) *storageFacade {
			return NewStorageFacade(storage, relyingParty, privateKey, staticClients...)
		}, fx.As(new(op.Storage)))),
		fx.Provide(func(storage op.Storage, configuration delegatedauth.Config) (op.OpenIDProvider, error) {
			keySet, err := ReadKeySet(context.TODO(), configuration)
			if err != nil {
				return nil, err
			}

			return NewOpenIDProvider(context.TODO(), storage, baseUrl.String(), configuration.Issuer, *keySet)
		}),
	)
}

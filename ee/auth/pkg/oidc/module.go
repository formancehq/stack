package oidc

import (
	"context"
	"crypto/rsa"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gopkg.in/square/go-jose.v2"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/delegatedauth"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/op"
	"go.uber.org/fx"
)

func Module(privateKey *rsa.PrivateKey, issuer string, staticClients ...auth.StaticClient) fx.Option {
	return fx.Options(
		fx.Invoke(fx.Annotate(func(router chi.Router, provider op.OpenIDProvider,
			storage Storage, relyingParty rp.RelyingParty) {
			AddRoutes(router, provider, storage, relyingParty)
		}, fx.ParamTags(``, ``, ``, `optional:"true"`))),
		fx.Provide(fx.Annotate(func(storage Storage, relyingParty rp.RelyingParty) *storageFacade {
			return NewStorageFacade(storage, relyingParty, privateKey, staticClients...)
		}, fx.As(new(op.Storage)), fx.ParamTags(``, `optional:"true"`))),
		fx.Provide(fx.Annotate(func(httpClient *http.Client, storage op.Storage, configuration delegatedauth.Config) (op.OpenIDProvider, error) {
			var (
				keySet *jose.JSONWebKeySet
				err    error
			)
			if configuration.Issuer != "" {
				keySet, err = ReadKeySet(httpClient, context.TODO(), configuration)
				if err != nil {
					return nil, err
				}
			}

			return NewOpenIDProvider(storage, issuer, configuration.Issuer, keySet)
		}, fx.ParamTags(``, ``, `optional:"true"`))),
	)
}

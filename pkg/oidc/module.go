package oidc

import (
	"context"
	"crypto/rsa"

	auth "github.com/formancehq/auth/pkg"
	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
	"go.uber.org/fx"
)

func Module(addr, issuer string, privateKey *rsa.PrivateKey) fx.Option {
	return fx.Options(
		fx.Provide(NewRouter),
		fx.Provide(fx.Annotate(func(storage Storage, relyingParty rp.RelyingParty, opts []auth.ClientOptions) *storageFacade {
			var staticClients []auth.Client
			for _, c := range opts {
				staticClients = append(staticClients, *auth.NewClient(c))
			}
			return NewStorageFacade(storage, relyingParty, privateKey, staticClients...)
		}, fx.As(new(op.Storage)))),
		fx.Provide(func(storage op.Storage) (op.OpenIDProvider, error) {
			return NewOpenIDProvider(context.TODO(), storage, issuer)
		}),
		fx.Invoke(func(lc fx.Lifecycle, router *mux.Router) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return StartServer(addr, router)
				},
			})
		}),
	)
}

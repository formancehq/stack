package delegatedauth

import (
	"context"

	"github.com/coreos/go-oidc"
)

type OIDCProvider = oidc.Provider

func ProvideDelegatedOIDCProvider(ctx context.Context, issuer Issuer) (*OIDCProvider, error) {
	return oidc.NewProvider(ctx, string(issuer))
}

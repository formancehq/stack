package api

import (
	"context"
	"crypto/sha256"

	"github.com/zitadel/oidc/pkg/op"
	"golang.org/x/text/language"
)

const (
	pathLoggedOut = "/logged-out"
)

type Issuer string

func NewOpenIDProvider(ctx context.Context, storage op.Storage, issuer Issuer) (op.OpenIDProvider, error) {
	return op.NewOpenIDProvider(ctx, &op.Config{
		Issuer:                   string(issuer),
		CryptoKey:                sha256.Sum256([]byte("test")),
		DefaultLogoutRedirectURI: pathLoggedOut,
		CodeMethodS256:           true,
		AuthMethodPost:           true,
		AuthMethodPrivateKeyJWT:  true,
		GrantTypeRefreshToken:    true,
		RequestObjectSupported:   true,
		SupportedUILocales:       []language.Tag{language.English},
	}, storage)
}

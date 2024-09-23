package oidc

import (
	"crypto/sha256"
	"net/http"
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
	"golang.org/x/text/language"
	"gopkg.in/go-jose/go-jose.v2"
)

const (
	pathLoggedOut = "/logged-out"
)

type verifier struct {
	issuer          string
	mat             time.Duration
	offset          time.Duration
	jsonWebKeySet   jose.JSONWebKeySet
	delegatedIssuer string
}

func (v verifier) DelegatedIssuer() string {
	return v.delegatedIssuer
}

func (v verifier) JSONWebKeySet() jose.JSONWebKeySet {
	return v.jsonWebKeySet
}

func (v verifier) Issuer() string {
	return v.issuer
}

func (v verifier) MaxAgeIAT() time.Duration {
	return v.mat
}

func (v verifier) Offset() time.Duration {
	return v.offset
}

type provider struct {
	op.OpenIDProvider
	delegatedIssuerJsonWebKeySet jose.JSONWebKeySet
	delegatedIssuer              string
	issuer                       string
}

func (p provider) JWTProfileVerifier() JWTProfileVerifier {
	return &verifier{
		issuer:          p.issuer,
		delegatedIssuer: p.delegatedIssuer,
		mat:             time.Hour,
		offset:          0,
		jsonWebKeySet:   p.delegatedIssuerJsonWebKeySet,
	}
}

var _ JWTAuthorizationGrantExchanger = (*provider)(nil)

func NewOpenIDProvider(storage op.Storage, issuer, delegatedIssuer string, delegatedIssuerJsonWebKeySet *jose.JSONWebKeySet) (op.OpenIDProvider, error) {
	var p op.OpenIDProvider

	interceptors := make([]op.Option, 0)
	if delegatedIssuer != "" {
		interceptors = append(interceptors, op.WithHttpInterceptors(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Intercept token requests with grant_type of type bearer assertion
				// as the library does not implement what we needs
				if r.URL.Path == op.DefaultEndpoints.Token.Relative() &&
					r.FormValue("grant_type") == string(oidc.GrantTypeBearer) {
					grantTypeBearer(issuer, &provider{
						issuer:                       issuer,
						OpenIDProvider:               p,
						delegatedIssuerJsonWebKeySet: *delegatedIssuerJsonWebKeySet,
						delegatedIssuer:              delegatedIssuer,
					}).ServeHTTP(w, r)
					return
				}
				handler.ServeHTTP(w, r)
			})

		}))
	}
	interceptors = append(interceptors, op.WithAllowInsecure())

	p, err := op.NewOpenIDProvider(issuer, &op.Config{
		CryptoKey:                sha256.Sum256([]byte("test")),
		DefaultLogoutRedirectURI: pathLoggedOut,
		CodeMethodS256:           true,
		AuthMethodPost:           true,
		AuthMethodPrivateKeyJWT:  true,
		GrantTypeRefreshToken:    true,
		RequestObjectSupported:   true,
		SupportedUILocales:       []language.Tag{language.English},
	}, storage, interceptors...)
	return p, err
}

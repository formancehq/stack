package oidc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v4"

	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type openIDKeySet struct {
	jose.JSONWebKeySet
}

// VerifySignature implements the oidc.KeySet interface
// providing an implementation for the keys stored in the OP Storage interface
func (o *openIDKeySet) VerifySignature(ctx context.Context, jws *jose.JSONWebSignature) ([]byte, error) {
	keyID, alg := oidc.GetKeyIDAndAlg(jws)

	var (
		key jose.JSONWebKey
		err error
	)
	key, err = oidc.FindMatchingKey(keyID, oidc.KeyUseSignature, alg, o.JSONWebKeySet.Keys...)
	if err != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}

	return jws.Verify(&key)
}

func VerifyJWTAssertion(ctx context.Context, assertion string, v JWTProfileVerifier) (*oidc.JWTTokenRequest, error) {
	request := new(oidc.JWTTokenRequest)

	_, err := oidc.ParseToken(assertion, request)
	if err != nil {
		return nil, err
	}

	if err := oidc.CheckAudience(request, v.Issuer()); err != nil {
		return nil, err
	}

	if err := oidc.CheckExpiration(request, v.Offset()); err != nil {
		return nil, err
	}

	accessTokenVerifier := op.NewAccessTokenVerifier(v.DelegatedIssuer(), &openIDKeySet{
		v.JSONWebKeySet(),
	})
	if _, err := op.VerifyAccessToken[*oidc.TokenClaims](ctx, assertion, accessTokenVerifier); err != nil {
		return nil, err
	}

	return request, nil
}

type JWTProfileVerifier interface {
	Issuer() string
	Offset() time.Duration
	DelegatedIssuer() string
	JSONWebKeySet() jose.JSONWebKeySet
}

type JWTAuthorizationGrantExchanger interface {
	op.Exchanger
	JWTProfileVerifier() JWTProfileVerifier
}

func grantTypeBearer(issuer string, p JWTAuthorizationGrantExchanger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileRequest, err := op.ParseJWTProfileGrantRequest(r, p.Decoder())
		if err != nil {
			op.RequestError(w, r, err, slog.Default())
			return
		}

		clientID, clientSecret, ok := r.BasicAuth()
		var client *clientFacade
		if ok {
			c, err := p.Storage().GetClientByClientID(r.Context(), clientID)
			if err != nil {
				op.RequestError(w, r, err, slog.Default())
				return
			}
			client = c.(*clientFacade)
			if !client.Client.IsPublic() {
				if err := client.Client.ValidateSecret(clientSecret); err != nil {
					op.RequestError(w, r, err, slog.Default())
					return
				}
			}
		}

		tokenRequest, err := VerifyJWTAssertion(r.Context(), profileRequest.Assertion, p.JWTProfileVerifier())
		if err != nil {
			op.RequestError(w, r, err, slog.Default())
			return
		}

		tokenRequest.Scopes, err = p.Storage().ValidateJWTProfileScopes(r.Context(), tokenRequest.Issuer, profileRequest.Scope)
		if err != nil {
			op.RequestError(w, r, err, slog.Default())
			return
		}

		tokens, err := ParseAssertion(profileRequest.Assertion)
		if err != nil {
			op.RequestError(w, r, err, slog.Default())
			return
		}

		tokenRequest.Scopes = tokens.Scopes

		resp, err := CreateJWTTokenResponse(r.Context(), issuer, tokenRequest, p, client)
		if err != nil {
			op.RequestError(w, r, err, slog.Default())
			return
		}

		httphelper.MarshalJSON(w, resp)
	}
}

func ParseAssertion(assertion string) (*oidc.AccessTokenClaims, error) {
	var claims = new(oidc.AccessTokenClaims)

	_, err := oidc.ParseToken(assertion, claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func CreateJWTTokenResponse(ctx context.Context, issuer string, tokenRequest *oidc.JWTTokenRequest, creator op.TokenCreator, client op.Client) (*oidc.AccessTokenResponse, error) {
	id, exp, err := creator.Storage().CreateAccessToken(ctx, tokenRequest)
	if err != nil {
		return nil, err
	}

	accessToken, err := op.CreateJWT(ctx, issuer, tokenRequest, exp, id, client, creator.Storage())
	if err != nil {
		return nil, err
	}

	return &oidc.AccessTokenResponse{
		AccessToken: accessToken,
		TokenType:   oidc.BearerToken,
		ExpiresIn:   uint64(exp.Sub(time.Now().UTC()).Seconds()),
	}, nil
}

package interceptors

import (
	"context"
	"strings"

	"github.com/formancehq/stack/components/stargate/internal/server/grpc/opentelemetry"
	"github.com/golang-jwt/jwt"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingMetadata = errors.New("missing metadata")
	ErrInvalidToken    = errors.New("invalid token")
	ErrFetchingJWKKeys = errors.New("error fetching JWK keys")
)

type AuthInterceptor struct {
	jwksURL         string
	metricsRegistry opentelemetry.MetricsRegistry
}

func NewAuthInterceptor(jwksURL string, metricsRegistry opentelemetry.MetricsRegistry) *AuthInterceptor {
	return &AuthInterceptor{
		jwksURL:         jwksURL,
		metricsRegistry: metricsRegistry,
	}
}

func (a *AuthInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			a.metricsRegistry.UnAuthenticatedCalls().Add(ctx, 1)
			return status.Errorf(codes.Unauthenticated, ErrMissingMetadata.Error())
		}

		if !a.valid(md["authorization"]) {
			a.metricsRegistry.UnAuthenticatedCalls().Add(ctx, 1)
			return status.Errorf(codes.Unauthenticated, ErrInvalidToken.Error())
		}

		return handler(srv, ss)
	}
}

func (a *AuthInterceptor) valid(authorization []string) bool {
	if len(authorization) == 0 {
		return false
	}

	token := strings.TrimPrefix(authorization[0], "Bearer ")
	t, err := jwt.Parse(token, a.getKey)
	if err != nil {
		// TODO(polo): should we panic ?
		return false
	}

	return t.Valid
}

func (a *AuthInterceptor) getKey(token *jwt.Token) (interface{}, error) {
	set, err := jwk.Fetch(context.Background(), a.jwksURL)
	if err != nil {
		return nil, err
	}

	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.Wrap(ErrFetchingJWKKeys, "expecting JWT header to have string kid")
	}

	key, flag := set.LookupKeyID(keyID)
	if !flag {
		return nil, errors.Wrapf(ErrFetchingJWKKeys, "unable to find key %q", keyID)
	}

	var pubkey interface{}
	err = key.Raw(&pubkey)
	if err != nil {
		return nil, errors.Wrapf(ErrFetchingJWKKeys, "unable to find key %q: %v", keyID, err)
	}

	return pubkey, nil
}

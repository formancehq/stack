package interceptors

import (
	"context"
	"net/http"
	"strings"

	"github.com/formancehq/stack/components/stargate/internal/server/grpc/opentelemetry"
	"github.com/golang-jwt/jwt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
	"github.com/zitadel/oidc/pkg/client"
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
	httpClient      *http.Client
	metricsRegistry opentelemetry.MetricsRegistry
}

func NewAuthInterceptor(jwksURL string, maxRetriesJWKSFetchting int, metricsRegistry opentelemetry.MetricsRegistry) *AuthInterceptor {
	return &AuthInterceptor{
		jwksURL:         jwksURL,
		httpClient:      newHttpClient(maxRetriesJWKSFetchting),
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

		method := info.FullMethod
		switch method {
		case "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo":
			return handler(srv, ss)
		default:
		}

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
	discoveryConfiguration, err := client.Discover(a.jwksURL, a.httpClient)
	if err != nil {
		return nil, err
	}

	set, err := jwk.Fetch(context.Background(), discoveryConfiguration.JwksURI)
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

func newHttpClient(maxRetries int) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = maxRetries
	return retryClient.StandardClient()
}

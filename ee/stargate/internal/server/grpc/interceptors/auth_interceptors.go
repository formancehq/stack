package interceptors

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/formancehq/stack/components/stargate/internal/server/grpc/metrics"
	"github.com/formancehq/stack/libs/go-libs/logging"
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

	IsLetter = regexp.MustCompile(`^[a-z]+$`).MatchString
)

type AuthInterceptor struct {
	logger          logging.Logger
	jwksURL         string
	httpClient      *http.Client
	metricsRegistry metrics.MetricsRegistry
}

func NewAuthInterceptor(
	logger logging.Logger,
	jwksURL string,
	maxRetriesJWKSFetchting int,
	metricsRegistry metrics.MetricsRegistry,
) *AuthInterceptor {
	return &AuthInterceptor{
		logger:          logger,
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
		a.logger.Errorf("error parsing token: %v", err)
		return false
	}

	if !t.Valid {
		return false
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		a.logger.Error("invalid claims format")
		return false
	}

	audience := ""
	switch v := claims["aud"].(type) {
	case string:
		audience = v
	case []interface{}:
		if len(v) == 0 {
			a.logger.Errorf("invalid type for audience %v", v)
			return false
		}
		audience, ok = v[0].(string)
		if !ok {
			a.logger.Errorf("invalid type for audience %v", v)
			return false
		}
	default:
		a.logger.Errorf("invalid type for audience %v", claims["aud"])
		return false
	}

	l := strings.Split(audience, "_")
	if len(l) != 3 {
		a.logger.Errorf("invalid audience format: %s", audience)
		return false
	}

	organizationID := l[1]
	if len(organizationID) != 12 || !IsLetter(organizationID) {
		a.logger.Errorf("invalid organization_id format: %s for audience: %s", organizationID, audience)
		return false
	}

	stackID := l[2]
	if len(stackID) != 4 || !IsLetter(stackID) {
		a.logger.Errorf("invalid stack_id format: %s for audience: %s", stackID, audience)
		return false
	}

	return true
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

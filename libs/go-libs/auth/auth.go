package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
	"go.uber.org/zap"
)

type jwtAuth struct {
	logger      logging.Logger
	httpClient  *http.Client
	verifiers   map[string]op.AccessTokenVerifier // issuer -> verifier
	checkScopes bool
	service     string
}

func newOtlpHttpClient(maxRetries int) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = maxRetries
	return retryClient.StandardClient()
}

func newJWTAuth(
	logger logging.Logger,
	readKeySetMaxRetries int,
	verifiers map[string]op.AccessTokenVerifier,
	service string,
	checkScopes bool,
) *jwtAuth {
	return &jwtAuth{
		logger:      logger,
		httpClient:  newOtlpHttpClient(readKeySetMaxRetries),
		verifiers:   verifiers,
		checkScopes: checkScopes,
		service:     service,
	}
}

// Authenticate validates the JWT in the request and returns the user, if valid.
func (ja *jwtAuth) Authenticate(w http.ResponseWriter, r *http.Request) (bool, error) {
	authHeader := r.Header.Get("authorization")
	if authHeader == "" {
		ja.logger.Error("no authorization header")
		return false, fmt.Errorf("no authorization header")
	}

	if !strings.HasPrefix(authHeader, strings.ToLower(oidc.PrefixBearer)) &&
		!strings.HasPrefix(authHeader, oidc.PrefixBearer) {
		ja.logger.Error("malformed authorization header")
		return false, fmt.Errorf("malformed authorization header")
	}

	token := strings.TrimPrefix(authHeader, strings.ToLower(oidc.PrefixBearer))
	token = strings.TrimPrefix(token, oidc.PrefixBearer)

	// Pre-parse the token to extract the issuer claim, so we can select
	// the correct verifier (each issuer has its own key set).
	var preClaims oidc.TokenClaims
	if _, err := oidc.ParseToken(token, &preClaims); err != nil {
		ja.logger.Error("unable to parse token", zap.Error(err))
		return false, fmt.Errorf("unable to parse token: %w", err)
	}

	verifier, ok := ja.verifiers[preClaims.Issuer]
	if !ok {
		issuers := make([]string, 0, len(ja.verifiers))
		for iss := range ja.verifiers {
			issuers = append(issuers, iss)
		}
		ja.logger.Error("untrusted issuer",
			zap.String("got", preClaims.Issuer),
			zap.Strings("trusted", issuers),
		)
		return false, fmt.Errorf("issuer does not match: got: %s, trusted: %v", preClaims.Issuer, issuers)
	}

	claims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](r.Context(), token, verifier)
	if err != nil {
		ja.logger.Error("unable to verify access token", zap.Error(err))
		return false, fmt.Errorf("unable to verify access token: %w", err)
	}

	if ja.checkScopes {
		scope := claims.Scopes

		allowed := true
		switch r.Method {
		case http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodTrace:
			allowed = allowed && (collectionutils.Contains(scope, ja.service+":read") || collectionutils.Contains(scope, ja.service+":write"))
		default:
			allowed = allowed && collectionutils.Contains(scope, ja.service+":write")
		}

		if !allowed {
			ja.logger.Info("not enough scopes")
			return false, fmt.Errorf("missing access, found scopes: '%s' need %s:read|write", strings.Join(scope, ", "), ja.service)
		}
	}

	return true, nil
}

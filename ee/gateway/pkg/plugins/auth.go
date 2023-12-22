package plugins

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/caddyauth"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/oidc/v2/pkg/op"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(JWTAuth{})
	httpcaddyfile.RegisterHandlerDirective("auth", parseAuthCaddyfile)
}

type JWTAuth struct {
	logger              *zap.Logger            `json:"-"`
	httpClient          *http.Client           `json:"-"`
	accessTokenVerifier op.AccessTokenVerifier `json:"-"`

	Issuer               string `json:"issuer,omitempty"`
	ReadKeySetMaxRetries int    `json:"read_key_set_max_retries,omitempty"`
	CheckScopes          bool   `json:"check_scopes,omitempty"`
	Service              string `json:"service"`
}

func (JWTAuth) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.authentication.providers.jwt",
		New: func() caddy.Module { return new(JWTAuth) },
	}
}

// Provision implements caddy.Provisioner interface.
func (ja *JWTAuth) Provision(ctx caddy.Context) error {
	ja.logger = ctx.Logger(ja)
	ja.httpClient = newOtlpHttpClient(ja.ReadKeySetMaxRetries)

	// We can't provision the access token verifier here because the auth
	// components is not started yet. He will not be started until the
	// gateway is started, which will be the case after this provision call.
	ja.accessTokenVerifier = nil
	return nil
}

func parseAuthCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var ja JWTAuth

	for h.Next() {
		for h.NextBlock(0) {
			opt := h.Val()
			switch opt {
			case "issuer":
				if !h.AllArgs(&ja.Issuer) {
					return nil, h.Errf("invalid issuer")
				}
			case "read_key_set_max_retries":
				var readKeySetMaxRetries string
				if !h.AllArgs(&readKeySetMaxRetries) {
					return nil, h.Errf("invalid read_key_set_max_retries")
				}

				var err error
				ja.ReadKeySetMaxRetries, err = strconv.Atoi(readKeySetMaxRetries)
				if err != nil {
					return nil, h.Errf("invalid read_key_set_max_retries")
				}
			case "check_scopes":
				var checkScopes string
				if !h.AllArgs(&checkScopes) {
					return nil, h.Errf("invalid check_scopes")
				}
				checkScopes = strings.ToLower(checkScopes)
				ja.CheckScopes = checkScopes == "true" || checkScopes == "1" || checkScopes == "yes"
			case "service":
				var service string
				if !h.AllArgs(&service) {
					return nil, h.Errf("invalid service")
				}
				ja.Service = service
			default:
				return nil, h.Errf("unrecognized option: %s", opt)
			}
		}
	}

	return caddyauth.Authentication{
		ProvidersRaw: caddy.ModuleMap{
			"jwt": caddyconfig.JSON(ja, nil),
		},
	}, nil
}

// Authenticate validates the JWT in the request and returns the user, if valid.
func (ja *JWTAuth) Authenticate(w http.ResponseWriter, r *http.Request) (caddyauth.User, bool, error) {
	authHeader := r.Header.Get("authorization")
	if authHeader == "" {
		ja.logger.Error("no authorization header")
		return caddyauth.User{}, false, fmt.Errorf("no authorization header")
	}

	if !strings.HasPrefix(authHeader, strings.ToLower(oidc.PrefixBearer)) &&
		!strings.HasPrefix(authHeader, oidc.PrefixBearer) {
		ja.logger.Error("malformed authorization header")
		return caddyauth.User{}, false, fmt.Errorf("malformed authorization header")
	}

	token := strings.TrimPrefix(authHeader, strings.ToLower(oidc.PrefixBearer))
	token = strings.TrimPrefix(token, oidc.PrefixBearer)

	accessTokenVerifier, err := ja.getAccessTokenVerifier(r.Context())
	if err != nil {
		ja.logger.Error("unable to create access token verifier", zap.Error(err))
		return caddyauth.User{}, false, fmt.Errorf("unable to create access token verifier: %w", err)
	}

	claims, err := op.VerifyAccessToken[*oidc.AccessTokenClaims](r.Context(), token, accessTokenVerifier)
	if err != nil {
		ja.logger.Error("unable to verify access token", zap.Error(err))
		return caddyauth.User{}, false, fmt.Errorf("unable to verify access token: %w", err)
	}

	if ja.CheckScopes {
		scope := claims.Scopes

		allowed := true
		switch r.Method {
		case http.MethodOptions, http.MethodGet, http.MethodHead, http.MethodTrace:
			allowed = allowed && (collectionutils.Contains(scope, ja.Service+":read") || collectionutils.Contains(scope, ja.Service+":write"))
		default:
			allowed = allowed && collectionutils.Contains(scope, ja.Service+":write")
		}

		if !allowed {
			ja.logger.Info("not enough scopes")
			return caddyauth.User{}, false, fmt.Errorf("missing access, found scopes: '%s' need %s:read|write", strings.Join(scope, ", "), ja.Service)
		}
	}

	return caddyauth.User{}, true, nil
}

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

func (ja *JWTAuth) getAccessTokenVerifier(ctx context.Context) (op.AccessTokenVerifier, error) {
	if ja.accessTokenVerifier == nil {
		//discoveryConfiguration, err := client.Discover(ja.Issuer, ja.httpClient)
		//if err != nil {
		//	return nil, err
		//}

		// todo: ugly quick fix
		authServicePort := "8080"
		if fromEnv := os.Getenv("AUTH_SERVICE_PORT"); fromEnv != "" {
			authServicePort = fromEnv
		}
		keySet := rp.NewRemoteKeySet(ja.httpClient, fmt.Sprintf("http://auth:%s/keys", authServicePort))

		ja.accessTokenVerifier = op.NewAccessTokenVerifier(
			os.Getenv("STACK_PUBLIC_URL")+"/api/auth",
			keySet,
		)
	}

	return ja.accessTokenVerifier, nil
}

func newOtlpHttpClient(maxRetries int) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = maxRetries
	return retryClient.StandardClient()
}

//------------------------------------------------------------------------------

// Interface Guards
var (
	_ caddy.Provisioner       = (*JWTAuth)(nil)
	_ caddy.Module            = (*JWTAuth)(nil)
	_ caddyauth.Authenticator = (*JWTAuth)(nil)
)

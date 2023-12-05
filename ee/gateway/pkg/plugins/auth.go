package plugins

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp/caddyauth"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zitadel/oidc/pkg/client"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/oidc/pkg/op"
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

	if _, err := op.VerifyAccessToken(r.Context(), token, accessTokenVerifier); err != nil {
		ja.logger.Error("unable to verify access token", zap.Error(err))
		return caddyauth.User{}, false, fmt.Errorf("unable to verify access token: %w", err)
	}

	return caddyauth.User{}, true, nil
}

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

func (ja *JWTAuth) getAccessTokenVerifier(
	ctx context.Context,
) (op.AccessTokenVerifier, error) {
	if ja.accessTokenVerifier == nil {
		discoveryConfiguration, err := client.Discover(ja.Issuer, ja.httpClient)
		if err != nil {
			return nil, err
		}

		keySet := rp.NewRemoteKeySet(ja.httpClient, discoveryConfiguration.JwksURI)

		ja.accessTokenVerifier = op.NewAccessTokenVerifier(
			ja.Issuer,
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

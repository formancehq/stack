package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func init() {
	caddy.RegisterModule(Versions{})
	httpcaddyfile.RegisterHandlerDirective("versions", parseCaddyfile)
}

type Endpoint struct {
	Name     string `json:"name,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

// Versions is a module that serves a /versions endpoint. This endpoint will
// gather all sub-services versions and return them in a JSON response.
// This module is configurable by end-users via caddy configuration.
type Versions struct {
	logger          *zap.Logger  `json:"-"`
	versionsHandler http.Handler `json:"-"`

	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

// Implements the caddy.Module interface.
func (Versions) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		// Note: The ID must start by the namespace http.handlers.* in order to be
		// loaded by the global handlers.
		ID:  "http.handlers.versions",
		New: func() caddy.Module { return new(Versions) },
	}
}

// Implements the caddy.Provisioner interface.
func (v *Versions) Provision(ctx caddy.Context) error {
	v.logger = ctx.Logger(v)
	v.versionsHandler = newVersionsHandler(
		v.logger,
		newHTTPClient(),
		v.Endpoints,
	)

	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var v Versions
	err := v.UnmarshalCaddyfile(h.Dispenser)
	return v, err
}

// UnmarshalCaddyfile sets up the handler from Caddyfile tokens. Syntax:
//
//	versions {
//		<field> <value>
//	}
func (m *Versions) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	m.Endpoints = make([]Endpoint, 0)
	for d.Next() {
		for d.NextBlock(0) {
			key := d.Val()
			var endpoint string
			d.Args(&endpoint)
			m.Endpoints = append(m.Endpoints, Endpoint{
				Name:     key,
				Endpoint: endpoint,
			})
		}
	}
	return nil
}

// Implements the caddy.Validator interface.
// Validate is called after the config initialization is completed and after
// the module is provisioned.
func (v *Versions) Validate() error {
	for _, endpoint := range v.Endpoints {
		if _, err := url.ParseRequestURI(endpoint.Endpoint); err != nil {
			return fmt.Errorf("invalid endpoint %s: %w", endpoint.Name, err)
		}
	}

	return nil
}

func (v Versions) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	v.versionsHandler.ServeHTTP(w, r)
	return nil
}

//------------------------------------------------------------------------------

type serviceInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type versionsResponse struct {
	Versions []serviceInfo `json:"versions"`
}

type versionsHandler struct {
	logger     *zap.Logger
	httpClient *http.Client

	endpoints []Endpoint
}

func newVersionsHandler(logger *zap.Logger, httpClient *http.Client, endpoints []Endpoint) http.Handler {
	return &versionsHandler{
		logger:     logger,
		httpClient: httpClient,
		endpoints:  endpoints,
	}
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func (v *versionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eg, ctxGroup := errgroup.WithContext(r.Context())

	versions := make(chan serviceInfo, len(v.endpoints))

	for _, endpoint := range v.endpoints {
		edpt := endpoint
		eg.Go(func() error {
			req, err := http.NewRequestWithContext(ctxGroup, http.MethodGet, edpt.Endpoint, http.NoBody)
			if err != nil {
				return fmt.Errorf("failed to create version request: %w", err)
			}

			res := serviceInfo{
				Name:    edpt.Name,
				Version: "unknown",
			}

			resp, err := v.httpClient.Do(req)
			if err != nil {
				v.logger.Error("failed to get version", zap.String("name", res.Name), zap.Error(err))
				versions <- res
				return nil
			}

			defer func() {
				err = resp.Body.Close()
				if err != nil {
					v.logger.Error("failed to close response body", zap.Error(err))
				}
			}()

			if resp.StatusCode != http.StatusOK {
				v.logger.Error("failed to get version", zap.String("name", res.Name), zap.String("status", resp.Status))
				versions <- res
				return nil
			}

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				v.logger.Error("failed to read response body", zap.String("name", res.Name), zap.Error(err))
				versions <- res
				return nil
			}

			if err := json.Unmarshal(responseBody, &res); err != nil {
				v.logger.Error("failed to unmarshal response body", zap.String("name", res.Name), zap.Error(err))
			}

			versions <- res

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		v.logger.Error("failed to query versions", zap.Error(err))
		return
	}

	close(versions)

	res := versionsResponse{}
	res.Versions = make([]serviceInfo, 0, len(v.endpoints))
	for version := range versions {
		res.Versions = append(res.Versions, version)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		v.logger.Error("failed to encode response", zap.Error(err))
	}
}

//------------------------------------------------------------------------------

// Interface Guards
var (
	_ caddy.Provisioner           = (*Versions)(nil)
	_ caddy.Module                = (*Versions)(nil)
	_ caddyhttp.MiddlewareHandler = (*Versions)(nil)
	_ caddyfile.Unmarshaler       = (*Versions)(nil)
	_ caddy.Validator             = (*Versions)(nil)
)

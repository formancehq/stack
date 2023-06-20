package plugins

import (
	"context"
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
	Name            string `json:"name,omitempty"`
	VersionEndpoint string `json:"endpoint,omitempty"`
	HealthEndpoint  string `json:"health,omitempty"`
}

// Versions is a module that serves a /versions endpoint. This endpoint will
// gather all sub-services versions and return them in a JSON response.
// This module is configurable by end-users via caddy configuration.
type Versions struct {
	logger          *zap.Logger  `json:"-"`
	versionsHandler http.Handler `json:"-"`

	Region      string     `json:"region,omitempty"`
	Environment string     `json:"env,omitempty"`
	Endpoints   []Endpoint `json:"endpoints,omitempty"`
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
		v.Region,
		v.Environment,
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
			switch key {
			case "endpoints":
				for d.Next() {
					for d.NextBlock(0) {
						name := d.Val()
						var versionEndpoint string
						var healthEndpoint string
						fmt.Println("remaining args for endpoints", d.CountRemainingArgs())
						if !d.AllArgs(&versionEndpoint, &healthEndpoint) {
							return d.Errf("invalid number of endpoints' arguments: want <name> <version_endpoint> <health_endpoint>")
						}
						m.Endpoints = append(m.Endpoints, Endpoint{
							Name:            name,
							VersionEndpoint: versionEndpoint,
							HealthEndpoint:  healthEndpoint,
						})
					}
				}

			case "region":
				if !d.AllArgs(&m.Region) {
					return d.Errf("invalid number of region's arguments: want <region>")
				}

			case "env":
				if !d.AllArgs(&m.Environment) {
					return d.Errf("invalid number of env's arguments: want <env>")
				}
			}
		}
	}
	return nil
}

// Implements the caddy.Validator interface.
// Validate is called after the config initialization is completed and after
// the module is provisioned.
func (v *Versions) Validate() error {
	for _, endpoint := range v.Endpoints {
		if _, err := url.ParseRequestURI(endpoint.VersionEndpoint); err != nil {
			return fmt.Errorf("invalid version endpoint %s: %w", endpoint.Name, err)
		}

		if _, err := url.ParseRequestURI(endpoint.HealthEndpoint); err != nil {
			return fmt.Errorf("invalid health endpoint %s: %w", endpoint.Name, err)
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
	Name string `json:"name"`
	// We do not want to omit empty values in the json response
	Version string `json:"version"`
	Health  bool   `json:"health"`
}

type backendServiceInfo struct {
	serviceInfo
	// Deprecated: ledger v1
	Data struct {
		Version string `json:"version"`
	}
}

func (info backendServiceInfo) GetVersion() string {
	if info.Data.Version != "" {
		return info.Data.Version
	}
	if info.Version != "" {
		return info.Version
	}
	return "unknown"
}

type versionsResponse struct {
	Region   string         `json:"region"`
	Env      string         `json:"env"`
	Versions []*serviceInfo `json:"versions"`
}

type versionsHandler struct {
	logger     *zap.Logger
	httpClient *http.Client

	region    string
	env       string
	endpoints []Endpoint
}

func newVersionsHandler(logger *zap.Logger, httpClient *http.Client, region, env string, endpoints []Endpoint) http.Handler {
	return &versionsHandler{
		logger:     logger,
		httpClient: httpClient,
		region:     region,
		env:        env,
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

	versions := make(chan *serviceInfo, len(v.endpoints))

	for _, endpoint := range v.endpoints {
		edpt := endpoint
		eg.Go(func() error {
			res := &serviceInfo{
				Name: edpt.Name,
			}

			version, err := serviceVersion(ctxGroup, v.httpClient, edpt.VersionEndpoint)
			if err != nil {
				// Log and Discard the error if there is any, and provide an
				// "unknown" version.
				v.logger.Error("failed to query version", zap.Error(err))
				res.Version = "unknown"
			} else {
				res.Version = version
			}

			health, err := serviceHealth(ctxGroup, v.httpClient, edpt.HealthEndpoint)
			if err != nil {
				// Log and Discard the error if there is any, and provide a
				// "false" health.
				v.logger.Error("failed to query health", zap.Error(err))
				res.Health = false
			} else {
				res.Health = health
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

	res := versionsResponse{
		Region: v.region,
		Env:    v.env,
	}
	res.Versions = make([]*serviceInfo, 0, len(v.endpoints))
	for version := range versions {
		res.Versions = append(res.Versions, version)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		v.logger.Error("failed to encode response", zap.Error(err))
	}
}

func serviceVersion(
	ctx context.Context,
	httpClient *http.Client,
	versionEndpoint string,
) (string, error) {
	sInfo := &backendServiceInfo{}

	resp, err := serviceCall(ctx, httpClient, versionEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to get version for %s: %w", versionEndpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get version for %s: %s", versionEndpoint, resp.Status)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body for %s: %w", versionEndpoint, err)
	}

	if err := json.Unmarshal(responseBody, &sInfo); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body for %s: %w", versionEndpoint, err)
	}

	return sInfo.GetVersion(), nil
}

func serviceHealth(
	ctx context.Context,
	httpClient *http.Client,
	healthEndpoint string,
) (bool, error) {
	resp, err := serviceCall(ctx, httpClient, healthEndpoint)
	if err != nil {
		return false, fmt.Errorf("failed to get health for %s: %w", healthEndpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("wrong status code for health endpoint %s: %d", healthEndpoint, resp.StatusCode)
	}

	return true, nil
}

func serviceCall(
	ctx context.Context,
	httpClient *http.Client,
	endpoint string,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create version request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get version for %s: %w", endpoint, err)
	}

	return resp, nil
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

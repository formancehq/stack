package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/caddyserver/caddy/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func init() {
	caddy.RegisterModule(version{})
}

type endpoint struct {
	Name     string `json:"name,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

type version struct {
	logger     *zap.Logger  `json:"-"`
	httpClient *http.Client `json:"-"`

	Endpoints []endpoint `json:"endpoints,omitempty"`
}

// Implements the caddy.Module interface.
func (version) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		// Note: The ID must start by the namespace admin.api.* in order to be
		// loaded by the admin API.
		ID:  "admin.api.versions",
		New: func() caddy.Module { return new(version) },
	}
}

// Implements the caddy.Provisioner interface.
func (v *version) Provision(ctx caddy.Context) error {
	v.logger = ctx.Logger(v)
	v.httpClient = newHTTPClient()

	return nil
}

// Implements the caddy.Validator interface.
// Validate is called after the config initialization is completed and after
// the module is provisioned.
func (v *version) Validate() error {
	for _, endpoint := range v.Endpoints {
		if _, err := url.ParseRequestURI(endpoint.Endpoint); err != nil {
			return fmt.Errorf("invalid endpoint %s: %w", endpoint.Name, err)
		}
	}

	return nil
}

// Implements the caddy.AdminRouter interface.
func (v *version) Routes() []caddy.AdminRoute {
	return []caddy.AdminRoute{
		{
			Pattern: "/versions",
			Handler: caddy.AdminHandlerFunc(v.handleVersions),
		},
	}
}

//------------------------------------------------------------------------------

type serviceInfo struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type versionsResponse struct {
	Versions []serviceInfo `json:"versions"`
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func (v *version) handleVersions(w http.ResponseWriter, r *http.Request) error {
	eg, ctxGroup := errgroup.WithContext(r.Context())

	versions := make(chan serviceInfo, len(v.Endpoints))

	for _, endpoint := range v.Endpoints {
		edpt := endpoint
		eg.Go(func() error {
			req, err := http.NewRequestWithContext(ctxGroup, http.MethodGet, edpt.Endpoint, http.NoBody)
			if err != nil {
				return fmt.Errorf("failed to create version request: %w", err)
			}

			resp, err := v.httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed to get version: %w", err)
			}

			defer func() {
				err = resp.Body.Close()
				if err != nil {
					v.logger.Error("failed to close response body", zap.Error(err))
				}
			}()

			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}

			res := serviceInfo{}
			if err = json.Unmarshal(responseBody, &res); err != nil {
				return fmt.Errorf("failed to unmarshal response body: %w", err)
			}

			res.Name = edpt.Name

			versions <- res

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	close(versions)

	res := versionsResponse{}
	res.Versions = make([]serviceInfo, 0, len(v.Endpoints))
	for version := range versions {
		res.Versions = append(res.Versions, version)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}

//------------------------------------------------------------------------------

// Interface Guards
var (
	_ caddy.Provisioner = (*version)(nil)
	_ caddy.Module      = (*version)(nil)
	_ caddy.AdminRouter = (*version)(nil)
	_ caddy.Validator   = (*version)(nil)
)

package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/onsi/ginkgo/v2"
)

type proxy struct {
	reverse *httputil.ReverseProxy
	url     *url.URL
}

var (
	gatewayServer *httptest.Server
	proxies       = map[string]proxy{}
)

func registerService(s string, url *url.URL) {
	proxies[s] = proxy{
		reverse: httputil.NewSingleHostReverseProxy(url),
		url:     url,
	}
}

type serviceInfo struct {
	Name string `json:"name"`
	// We do not want to omit empty values in the json response
	Version string `json:"version"`
	Health  bool   `json:"health"`
}

type versionsResponse struct {
	Versions []*serviceInfo `json:"versions"`
}

func startFakeGateway() {
	gatewayServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/versions" {
			res := versionsResponse{
				Versions: []*serviceInfo{
					{
						Name:    "ledger",
						Version: "v2.0.0",
					},
					// If needed, add other services version
				},
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				panic(err)
			}
			return
		}
		for name, proxy := range proxies {
			if strings.HasPrefix(r.URL.Path, "/api/"+name) {
				ginkgo.GinkgoWriter.Printf("Proxying %s: %s\r\n", name, proxy.url.String())
				http.StripPrefix("/api/"+name, proxy.reverse).ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
}

func stopFakeGateway() {
	gatewayServer.Close()
}

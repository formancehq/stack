package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type serviceInfo struct {
	Name string `json:"name"`
	// We do not want to omit empty values in the json response
	Version string `json:"version"`
	Health  bool   `json:"health"`
}

type versionsResponse struct {
	Versions []*serviceInfo `json:"versions"`
}

type gateway struct {
	test *Test
}

func (g gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if !strings.HasPrefix(r.URL.Path, "/api/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	paths := strings.Split(r.URL.Path, "/")
	service := paths[2]
	routings, ok := g.test.servicesToRoute[service]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	found := false
	var port uint16
	for _, routing := range routings {
		if routing.routingFunc == nil {
			port = routing.port
			found = true
			break
		}

		ok = routing.routingFunc("/"+strings.Join(paths[3:], "/"), r.Method)
		if ok {
			port = routing.port
			found = true
			break
		}
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))
	proxy := httputil.NewSingleHostReverseProxy(url)

	http.StripPrefix("/api/"+service, proxy).ServeHTTP(w, r)
}

var _ http.Handler = &gateway{}

func newGateway(test *Test) *gateway {
	return &gateway{
		test: test,
	}
}

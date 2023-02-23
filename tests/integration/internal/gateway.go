package internal

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
)

var (
	gatewayServer *httptest.Server
	proxies       = map[string]*httputil.ReverseProxy{}
)

func registerService(s string, url *url.URL) {
	proxies[s] = httputil.NewSingleHostReverseProxy(url)
}

func startFakeGateway() {
	gatewayServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for name, proxy := range proxies {
			if strings.HasPrefix(r.URL.Path, "/api/"+name) {
				http.StripPrefix("/api/"+name, proxy).ServeHTTP(w, r)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func stopFakeGateway() {
	gatewayServer.Close()
}

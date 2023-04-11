package internal

import (
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

func startFakeGateway() {
	gatewayServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

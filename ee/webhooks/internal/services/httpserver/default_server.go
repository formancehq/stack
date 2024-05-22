package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/internal/commons"
	"github.com/formancehq/webhooks/internal/services/httpserver/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
)


var HealthcheckRoute = "/_healthcheck"
var InfoRoute = "/_info"


type DefaultServerParams struct {
	Addr string
	Info commons.ServiceInfo
	Auth auth.Auth
	Logger logging.Logger
}


type DefaultHTTPServer struct {
	server *http.Server
	handler *chi.Mux
	logger logging.Logger
}

func (server *DefaultHTTPServer) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", server.server.Addr)
	if err != nil {
		return err
	}
	logging.FromContext(ctx).Infof("Start Webhook HTTP server at %s", server.server.Addr)
	go server.server.Serve(ln)
	return nil
}

func (server *DefaultHTTPServer) Stop(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}

func (server *DefaultHTTPServer) Register(method string, url string, handler func(http.ResponseWriter, *http.Request)){
	switch MethodHTTP(method) {
	case GET :
		server.handler.Get(url, handler)
	case POST : 
		server.handler.Post(url, handler)
	case PUT :
		server.handler.Put(url, handler)
	case DELETE :
		server.handler.Delete(url, handler) 
	}
}

func NewDefaultHTTPServer(addr string, info commons.ServiceInfo, a auth.Auth, logger logging.Logger ) DefaultHTTPServer{
	
	router := chi.NewRouter()

	defaultServer := DefaultHTTPServer{	
		handler : router,
		server :  &http.Server{
			Addr: addr,
			Handler: router,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}

	defaultServer.handler.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	defaultServer.handler.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r.WithContext(logging.ContextWithLogger(r.Context(), logger)))
		})
	})
	defaultServer.handler.Group(func(r chi.Router) {
		r.Get(HealthcheckRoute, func(_ http.ResponseWriter, r *http.Request) {
			logging.FromContext(r.Context()).Infof("health check OK")
		})
		r.Get(InfoRoute,  func(w http.ResponseWriter, r *http.Request) {
			sharedapi.RawOk(w, info)
		})	
	})
	defaultServer.handler.Use(auth.Middleware(a))
	defaultServer.handler.Use(otelchi.Middleware("webhooks"))
	
	return defaultServer
}

var _ interfaces.IHTTPServer = (*DefaultHTTPServer)(nil)
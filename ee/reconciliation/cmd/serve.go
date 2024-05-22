package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/reconciliation/internal/client"
	apiv1 "github.com/formancehq/reconciliation/internal/v1/api"
	backendv1 "github.com/formancehq/reconciliation/internal/v1/api/backend"
	storagev1 "github.com/formancehq/reconciliation/internal/v1/storage"
	apiv2 "github.com/formancehq/reconciliation/internal/v2/api"
	backendv2 "github.com/formancehq/reconciliation/internal/v2/api/backend"
	storagev2 "github.com/formancehq/reconciliation/internal/v2/storage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func stackClientModule() fx.Option {
	return fx.Options(
		fx.Provide(func() *sdk.Formance {
			oauthConfig := clientcredentials.Config{
				ClientID:     viper.GetString(stackClientIDFlag),
				ClientSecret: viper.GetString(stackClientSecretFlag),
				TokenURL:     fmt.Sprintf("%s/api/auth/oauth/token", viper.GetString(stackURLFlag)),
				Scopes:       []string{"openid", "ledger:read", "ledger:write", "payments:read", "payments:write"},
			}
			underlyingHTTPClient := &http.Client{
				Transport: otlp.NewRoundTripper(http.DefaultTransport, viper.GetBool(service.DebugFlag)),
			}
			return sdk.New(
				sdk.WithClient(
					oauthConfig.Client(context.WithValue(context.Background(),
						oauth2.HTTPClient, underlyingHTTPClient)),
				),
				sdk.WithServerURL(viper.GetString(stackURLFlag)),
			)
		}),
	)
}

func newServeCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "serve",
		RunE: runServer(version),
	}
	cmd.Flags().String(listenFlag, ":8080", "Listening address")
	cmd.Flags().Bool(workerFlag, false, "Enable worker mode")

	return cmd
}

func runServer(version string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		commonOptions, err := commonOptions(cmd)
		if err != nil {
			return err
		}

		options := []fx.Option{
			commonOptions,
		}

		options = append(options,
			stackClientModule(),
			serverOptions(sharedapi.ServiceInfo{
				Version: version,
			}, viper.GetString(listenFlag)),
		)

		if viper.GetBool(workerFlag) {
			options = append(options, workerOptions())
		}

		return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
	}
}

func prepareDatabaseOptions(cmd *cobra.Command) (fx.Option, error) {
	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(cmd.Context())
	if err != nil {
		return nil, err
	}

	options := make([]fx.Option, 0)

	options = append(options,
		storagev1.Module(*connectionOptions),
		storagev2.Module(*connectionOptions),
	)

	return fx.Options(options...), nil
}

func healthCheckModule() fx.Option {
	return fx.Options(
		health.Module(),
		health.ProvideHealthCheck(func() health.NamedCheck {
			return health.NewNamedCheck("default", health.CheckFn(func(ctx context.Context) error {
				return nil
			}))
		}),
	)
}

func serverOptions(serviceInfo sharedapi.ServiceInfo, bind string) fx.Option {
	return fx.Options(
		healthCheckModule(),
		client.Module(),
		fx.Supply(serviceInfo),
		fx.Invoke(func(m *chi.Mux, lc fx.Lifecycle) {
			lc.Append(httpserver.NewHook(m, httpserver.WithAddress(bind)))
		}),
		apiv1.HTTPModule(),
		apiv2.HTTPModule(),
		fx.Provide(newRouter),
	)
}

func newRouter(
	backendV1 backendv1.Backend,
	backendV2 backendv2.Backend,
	serviceInfo api.ServiceInfo,
	healthController *health.HealthController,
	a auth.Auth,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	r.Get("/_healthcheck", healthController.Check)
	r.Get("/_info", api.InfoHandler(serviceInfo))

	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(a))
		r.Use(otelchi.Middleware("reconciliation"))

		apiv1.NewRouter(backendV1, r)
		apiv2.NewRouter(backendV2, r)
	})

	return r
}

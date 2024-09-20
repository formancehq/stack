package cmd

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/opensearch-project/opensearch-go"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/aws/iam"
	"github.com/formancehq/go-libs/health"
	"github.com/formancehq/go-libs/httpclient"
	"github.com/formancehq/go-libs/httpserver"
	"github.com/formancehq/go-libs/licence"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/otlp"
	"github.com/formancehq/go-libs/otlp/otlptraces"
	"github.com/formancehq/go-libs/service"
	"github.com/formancehq/search/pkg/searchengine"
	"github.com/formancehq/search/pkg/searchhttp"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opensearch-project/opensearch-go/opensearchtransport"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const (
	serviceName = "search"

	openSearchServiceFlag    = "open-search-service"
	openSearchSchemeFlag     = "open-search-scheme"
	openSearchUsernameFlag   = "open-search-username"
	openSearchPasswordFlag   = "open-search-password"
	esIndicesFlag            = "es-indices"
	esDisableMappingInitFlag = "mapping-init-disabled"
	bindFlag                 = "bind"
	stackFlag                = "stack"
	awsIAMEnabledFlag        = "aws-iam-enabled"

	defaultBind = ":8080"

	healthCheckPath = "/_healthcheck"
)

func NewServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "serve",
		Short:        "Launch the search server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			esIndex, _ := cmd.Flags().GetString(esIndicesFlag)
			if esIndex == "" {
				return errors.New("es index not defined")
			}

			bind, _ := cmd.Flags().GetString(bindFlag)
			if bind == "" {
				bind = defaultBind
			}

			options := make([]fx.Option, 0)
			options = append(options, opensearchClientModule(cmd))
			options = append(options,
				auth.FXModuleFromFlags(cmd),
				health.Module(),
				health.ProvideHealthCheck(func(client *opensearch.Client) health.NamedCheck {
					return health.NewNamedCheck("elasticsearch connection", health.CheckFn(func(ctx context.Context) error {
						_, err := client.Ping()
						return err
					}))
				}),
			)

			options = append(options, otlptraces.FXModuleFromFlags(cmd))
			options = append(options, licence.FXModuleFromFlags(cmd, serviceName))
			options = append(options, apiModule(cmd))

			return service.New(cmd.OutOrStdout(), options...).Run(cmd)
		},
	}

	cmd.Flags().String(bindFlag, defaultBind, "http server address")
	cmd.Flags().String(esIndicesFlag, "", "ES index to look")
	cmd.Flags().String(openSearchServiceFlag, "", "Open search service hostname")
	cmd.Flags().String(openSearchSchemeFlag, "https", "OpenSearch scheme")
	cmd.Flags().String(openSearchUsernameFlag, "", "OpenSearch username")
	cmd.Flags().String(openSearchPasswordFlag, "", "OpenSearch password")
	cmd.Flags().Bool(esDisableMappingInitFlag, false, "Disable mapping initialization")
	cmd.Flags().String(stackFlag, "", "Stack id")
	cmd.Flags().Bool(awsIAMEnabledFlag, false, "Use AWS IAM for authentication")

	iam.AddFlags(cmd.Flags())
	otlptraces.AddFlags(cmd.Flags())
	service.AddFlags(cmd.Flags())
	licence.AddFlags(cmd.Flags())
	auth.AddFlags(cmd.Flags())

	return cmd
}

func exitWithError(ctx context.Context, msg string) {
	logging.FromContext(ctx).Error(msg)
	os.Exit(1)
}

func newOpensearchClient(cmd *cobra.Command, config *opensearch.Config) (*opensearch.Client, error) {
	httpTransport := http.DefaultTransport
	httpTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	httpTransport = otlp.NewRoundTripper(httpTransport, service.IsDebug(cmd))

	if service.IsDebug(cmd) {
		httpTransport = httpclient.NewDebugHTTPTransport(httpTransport)
		config.Logger = &opensearchtransport.JSONLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}
	} else {
		config.UseResponseCheckOnly = true
	}

	return opensearch.NewClient(*config)
}

func newConfig(cmd *cobra.Command) (*opensearch.Config, error) {
	openSearchServiceHost, _ := cmd.Flags().GetString(openSearchServiceFlag)
	if openSearchServiceHost == "" {
		exitWithError(cmd.Context(), "missing open search service host")
	}
	openSearchScheme, _ := cmd.Flags().GetString(openSearchSchemeFlag)
	awsIAMEnabled, _ := cmd.Flags().GetBool(awsIAMEnabledFlag)

	cfg := opensearch.Config{
		Addresses: []string{openSearchScheme + "://" + openSearchServiceHost},
	}
	if awsIAMEnabled {
		awsConfig, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, err
		}
		signer, err := requestsigner.NewSigner(awsConfig)
		if err != nil {
			return nil, err
		}
		cfg.Signer = signer
	} else {
		cfg.Username, _ = cmd.Flags().GetString(openSearchUsernameFlag)
		cfg.Password, _ = cmd.Flags().GetString(openSearchPasswordFlag)
	}
	return &cfg, nil
}

func opensearchClientModule(cmd *cobra.Command) fx.Option {
	esDisableMappingInit, _ := cmd.Flags().GetBool(esDisableMappingInitFlag)
	openSearchServiceHost, _ := cmd.Flags().GetString(openSearchServiceFlag)
	if openSearchServiceHost == "" {
		exitWithError(cmd.Context(), "missing open search service host")
	}
	esIndex, _ := cmd.Flags().GetString(esIndicesFlag)

	options := []fx.Option{
		fx.Provide(func() (*opensearch.Config, error) {
			return newConfig(cmd)
		}),
		fx.Provide(func(config *opensearch.Config) (*opensearch.Client, error) {
			return newOpensearchClient(cmd, config)
		}),
	}
	if !esDisableMappingInit {
		options = append(options, fx.Invoke(func(lc fx.Lifecycle, client *opensearch.Client) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return searchengine.CreateIndex(ctx, client, esIndex)
				},
			})
		}))
	}
	return fx.Options(options...)
}

func apiModule(cmd *cobra.Command) fx.Option {

	serviceInfo := api.ServiceInfo{
		Version: Version,
		Debug:   service.IsDebug(cmd),
	}
	stack, _ := cmd.Flags().GetString(stackFlag)
	esIndex, _ := cmd.Flags().GetString(esIndicesFlag)
	bind, _ := cmd.Flags().GetString(bindFlag)
	otelTraces, _ := cmd.Flags().GetBool(otlptraces.OtelTracesFlag)

	return fx.Options(
		fx.Provide(fx.Annotate(func(openSearchClient *opensearch.Client, tp trace.TracerProvider, healthController *health.HealthController, a auth.Authenticator) (http.Handler, error) {
			router := mux.NewRouter()

			router.Use(handlers.RecoveryHandler())
			router.Handle(healthCheckPath, http.HandlerFunc(healthController.Check))
			router.Path("/_info").Methods(http.MethodGet).Handler(api.InfoHandler(serviceInfo))

			routerWithTraces := router.PathPrefix("/").Subrouter()
			routerWithTraces.Use(auth.Middleware(a))
			if otelTraces {
				routerWithTraces.Use(otelmux.Middleware(serviceName, otelmux.WithTracerProvider(tp)))
			}
			routerWithTraces.PathPrefix("/").Handler(searchhttp.Handler(searchengine.NewDefaultEngine(
				openSearchClient,
				stack,
				searchengine.WithESIndex(esIndex),
			)))

			return router, nil
		}, fx.ParamTags(``, `optional:"true"`))),
		fx.Invoke(func(lc fx.Lifecycle, handler http.Handler) {
			lc.Append(httpserver.NewHook(handler, httpserver.WithAddress(bind)))
		}),
	)
}

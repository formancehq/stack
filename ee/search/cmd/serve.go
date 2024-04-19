package cmd

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/formancehq/search/pkg/searchengine"
	"github.com/formancehq/search/pkg/searchhttp"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpclient"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	app "github.com/formancehq/stack/libs/go-libs/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchtransport"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			openSearchServiceHost := viper.GetString(openSearchServiceFlag)
			if openSearchServiceHost == "" {
				exitWithError(cmd.Context(), "missing open search service host")
			}

			esIndex := viper.GetString(esIndicesFlag)
			if esIndex == "" {
				return errors.New("es index not defined")
			}

			bind := viper.GetString(bindFlag)
			if bind == "" {
				bind = defaultBind
			}

			options := make([]fx.Option, 0)
			options = append(options, opensearchClientModule(openSearchServiceHost, !viper.GetBool(esDisableMappingInitFlag), esIndex))
			options = append(options,
				auth.CLIAuthModule(),
				health.Module(),
				health.ProvideHealthCheck(func(client *opensearch.Client) health.NamedCheck {
					return health.NewNamedCheck("elasticsearch connection", health.CheckFn(func(ctx context.Context) error {
						_, err := client.Ping()
						return err
					}))
				}),
			)

			options = append(options, otlptraces.CLITracesModule())
			options = append(options, licence.CLIModule(serviceName))
			options = append(options, apiModule(serviceName, bind, viper.GetString(stackFlag), api.ServiceInfo{
				Version: Version,
			}, esIndex))

			return app.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
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
	iam.InitFlags(cmd.Flags())
	otlptraces.InitOTLPTracesFlags(cmd.Flags())
	app.BindFlags(cmd)
	licence.InitCLIFlags(cmd)

	return cmd
}

func exitWithError(ctx context.Context, msg string) {
	logging.FromContext(ctx).Error(msg)
	os.Exit(1)
}

func newOpensearchClient(config *opensearch.Config) (*opensearch.Client, error) {
	httpTransport := http.DefaultTransport
	httpTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	httpTransport = otlp.NewRoundTripper(httpTransport, viper.GetBool(app.DebugFlag))

	if viper.GetBool(app.DebugFlag) {
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

func newConfig(openSearchServiceHost string) (*opensearch.Config, error) {
	cfg := opensearch.Config{
		Addresses: []string{viper.GetString(openSearchSchemeFlag) + "://" + openSearchServiceHost},
	}
	if viper.GetBool(awsIAMEnabledFlag) {
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
		cfg.Username = viper.GetString(openSearchUsernameFlag)
		cfg.Password = viper.GetString(openSearchPasswordFlag)
	}
	return &cfg, nil
}

func opensearchClientModule(openSearchServiceHost string, loadMapping bool, esIndex string) fx.Option {

	options := []fx.Option{
		fx.Provide(func() (*opensearch.Config, error) {
			return newConfig(openSearchServiceHost)
		}),
		fx.Provide(func(config *opensearch.Config) (*opensearch.Client, error) {
			return newOpensearchClient(config)
		}),
	}
	if loadMapping {
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

func apiModule(
	serviceName, bind, stack string,
	serviceInfo api.ServiceInfo,
	esIndex string,
) fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(func(openSearchClient *opensearch.Client, tp trace.TracerProvider, healthController *health.HealthController, a auth.Auth) (http.Handler, error) {
			router := mux.NewRouter()

			router.Use(handlers.RecoveryHandler())
			router.Handle(healthCheckPath, http.HandlerFunc(healthController.Check))
			router.Path("/_info").Methods(http.MethodGet).Handler(api.InfoHandler(serviceInfo))

			routerWithTraces := router.PathPrefix("/").Subrouter()
			routerWithTraces.Use(auth.Middleware(a))
			if viper.GetBool(otlptraces.OtelTracesFlag) {
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

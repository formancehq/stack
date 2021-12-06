package main

import (
	"context"
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/numary/go-libs/oauth2/oauth2introspect"
	"github.com/numary/go-libs/sharedauth"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/go-libs/sharedotlp/sharedotlptraces"
	"github.com/numary/search/pkg/searchengine"
	"github.com/numary/search/pkg/searchhttp"
	"github.com/opensearch-project/opensearch-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gopkg.in/yaml.v3"
)

//go:embed swagger.yaml
var res embed.FS

const (
	openSearchServiceFlag                = "open-search-service"
	openSearchSchemeFlag                 = "open-search-scheme"
	esIndicesFlag                        = "es-indices"
	bindFlag                             = "bind"
	debugFlag                            = "debug"
	versionFlag                          = "version"
	otelTracesFlag                       = "otel-traces"
	otelTracesBatchFlag                  = "otel-traces-batch"
	otelTracesExporterFlag               = "otel-traces-exporter"
	otelTracesExporterJaegerEndpointFlag = "otel-traces-exporter-jaeger-endpoint"
	otelTracesExporterJaegerUserFlag     = "otel-traces-exporter-jaeger-user"
	otelTracesExporterJaegerPasswordFlag = "otel-traces-exporter-jaeger-password"
	otelTracesExporterOTLPModeFlag       = "otel-traces-exporter-otlp-mode"
	otelTracesExporterOTLPEndpointFlag   = "otel-traces-exporter-otlp-endpoint"
	otelTracesExporterOTLPInsecureFlag   = "otel-traces-exporter-otlp-insecure"

	authBasicEnabledFlag        = "auth-basic-enabled"
	authBasicCredentialsFlag    = "auth-basic-credentials"
	authBearerEnabledFlag       = "auth-bearer-enabled"
	authBearerIntrospectUrlFlag = "auth-bearer-introspect-url"
	authBearerAudienceFlag      = "auth-bearer-audience"

	defaultBind = ":8080"

	healthCheckPath = "/_healthcheck"
)

var Version = "latest"

func exitWithError(logger *logrus.Logger, msg string) {
	logger.Error(msg)
	os.Exit(1)
}

func yamlToJson(r io.Reader) string {
	m := make(map[string]interface{})
	dec := yaml.NewDecoder(r)
	err := dec.Decode(&m)
	if err != nil {
		panic(err)
	}
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func readSwaggerFile() string {
	file, err := res.Open("swagger.yaml")
	if err != nil {
		panic(err)
	}
	return yamlToJson(file)
}

func opensearchClientModule(openSearchServiceHost string, esIndices ...string) fx.Option {
	return fx.Options(
		fx.Provide(func() (*opensearch.Client, error) {
			httpTransport := http.DefaultTransport
			httpTransport.(*http.Transport).TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}

			return opensearch.NewClient(opensearch.Config{
				Addresses: []string{viper.GetString(openSearchSchemeFlag) + "://" + openSearchServiceHost},
				Transport: httpTransport,
			})
		}),
		fx.Invoke(func(lc fx.Lifecycle, client *opensearch.Client) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return searchengine.LoadDefaultMapping(context.TODO(), client, esIndices...)
				},
			})
		}),
	)
}

func apiModule(serviceName, bind string, esIndices ...string) fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(func(openSearchClient *opensearch.Client, tp trace.TracerProvider) (http.Handler, error) {
			swaggerFile := readSwaggerFile()
			router := mux.NewRouter()
			if viper.GetBool(otelTracesFlag) {
				router.Use(otelmux.Middleware(serviceName, otelmux.WithTracerProvider(tp)))
			}
			router.Use(handlers.RecoveryHandler())
			router.Handle("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(swaggerFile))
			}))
			router.Handle(healthCheckPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			protected := router.PathPrefix("/").Subrouter()

			methods := make([]sharedauth.Method, 0)
			if viper.GetBool(authBasicEnabledFlag) {
				credentials := sharedauth.Credentials{}
				for _, kv := range viper.GetStringSlice(authBasicCredentialsFlag) {
					parts := strings.SplitN(kv, ":", 2)
					credentials[parts[0]] = sharedauth.Credential{
						Password: parts[1],
						Scopes:   []string{"search"},
					}
				}
				methods = append(methods, sharedauth.NewHTTPBasicMethod(credentials))
			}
			if viper.GetBool(authBearerEnabledFlag) {
				methods = append(methods, sharedauth.NewHttpBearerMethod(
					sharedauth.NewIntrospectionValidator(
						oauth2introspect.NewIntrospecter(viper.GetString(authBearerIntrospectUrlFlag)),
						false,
						sharedauth.AudienceIn(viper.GetString(authBearerAudienceFlag)),
					),
				))
			}

			if len(methods) > 0 {
				protected.Use(sharedauth.Middleware(methods...))
			}
			router.PathPrefix("/").Handler(searchhttp.Handler(searchengine.NewDefaultEngine(
				openSearchClient,
				searchengine.WithESIndices(esIndices...),
			)))

			return router, nil
		}, fx.ParamTags(``, `optional:"true"`))),
		fx.Invoke(func(lc fx.Lifecycle, handler http.Handler) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					sharedlogging.GetLogger(ctx).Infof("Starting http server on %s", bind)
					go func() {
						err := http.ListenAndServe(bind, handler)
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
							os.Exit(1)
						}
					}()
					return nil
				},
			})
		}),
	)
}

func telemetryModule() fx.Option {
	return sharedotlptraces.TracesModule(sharedotlptraces.ModuleConfig{
		Batch:    viper.GetBool(otelTracesBatchFlag),
		Exporter: viper.GetString(otelTracesExporterFlag),
		JaegerConfig: func() *sharedotlptraces.JaegerConfig {
			if viper.GetString(otelTracesExporterFlag) != sharedotlptraces.JaegerExporter {
				return nil
			}
			return &sharedotlptraces.JaegerConfig{
				Endpoint: viper.GetString(otelTracesExporterJaegerEndpointFlag),
				User:     viper.GetString(otelTracesExporterJaegerUserFlag),
				Password: viper.GetString(otelTracesExporterJaegerPasswordFlag),
			}
		}(),
		OTLPConfig: func() *sharedotlptraces.OTLPConfig {
			if viper.GetString(otelTracesExporterFlag) != sharedotlptraces.OTLPExporter {
				return nil
			}
			return &sharedotlptraces.OTLPConfig{
				Mode:     viper.GetString(otelTracesExporterOTLPModeFlag),
				Endpoint: viper.GetString(otelTracesExporterOTLPEndpointFlag),
				Insecure: viper.GetBool(otelTracesExporterOTLPInsecureFlag),
			}
		}(),
	})
}

var cmd = cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {

		if viper.GetBool(versionFlag) {
			fmt.Println("Version:", Version)
			os.Exit(0)
		}
		debug := viper.GetBool(debugFlag)

		logger := logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{})
		if debug {
			logger.SetLevel(logrus.DebugLevel)
			logger.Debugln("Debug mode enabled")
		}
		logger.Debugf("Starting with config: %s", viper.AllSettings())

		openSearchServiceHost := viper.GetString(openSearchServiceFlag)
		if openSearchServiceHost == "" {
			exitWithError(logger, "missing open search service host")
		}

		esIndices := viper.GetStringSlice(esIndicesFlag)
		if len(esIndices) == 0 {
			esIndices = searchengine.DefaultEsIndices
		}

		bind := viper.GetString(bindFlag)
		if bind == "" {
			bind = defaultBind
		}

		options := make([]fx.Option, 0)
		options = append(options, opensearchClientModule(openSearchServiceHost, esIndices...))

		if viper.GetBool(otelTracesFlag) {
			options = append(options, telemetryModule())
		}
		options = append(options, apiModule("search", bind, esIndices...))

		app := fx.New(options...)

		err := app.Start(cmd.Context())
		if err != nil {
			return err
		}

		<-app.Done()

		return nil
	},
}

func init() {
	cmd.Flags().String(bindFlag, defaultBind, "http server address")
	cmd.Flags().StringSlice(esIndicesFlag, searchengine.DefaultEsIndices, "ES indices to look")
	cmd.Flags().String(openSearchServiceFlag, "", "Open search service hostname")
	cmd.Flags().String(openSearchSchemeFlag, "https", "OpenSearch scheme")
	cmd.Flags().Bool(otelTracesFlag, false, "Enable OpenTelemetry traces support")
	cmd.Flags().Bool(otelTracesBatchFlag, false, "Use OpenTelemetry batching")
	cmd.Flags().String(otelTracesExporterFlag, "stdout", "OpenTelemetry traces exporter")
	cmd.Flags().String(otelTracesExporterJaegerEndpointFlag, "", "OpenTelemetry traces Jaeger exporter endpoint")
	cmd.Flags().String(otelTracesExporterJaegerUserFlag, "", "OpenTelemetry traces Jaeger exporter user")
	cmd.Flags().String(otelTracesExporterJaegerPasswordFlag, "", "OpenTelemetry traces Jaeger exporter password")
	cmd.Flags().String(otelTracesExporterOTLPModeFlag, "grpc", "OpenTelemetry traces OTLP exporter mode (grpc|http)")
	cmd.Flags().String(otelTracesExporterOTLPEndpointFlag, "", "OpenTelemetry traces grpc endpoint")
	cmd.Flags().Bool(otelTracesExporterOTLPInsecureFlag, false, "OpenTelemetry traces grpc insecure")
	cmd.Flags().Bool(authBasicEnabledFlag, false, "Enable basic auth")
	cmd.Flags().StringSlice(authBasicCredentialsFlag, []string{}, "HTTP basic auth credentials (<username>:<password>)")
	cmd.Flags().Bool(authBearerEnabledFlag, false, "Enable bearer auth")
	cmd.Flags().String(authBearerIntrospectUrlFlag, "", "OAuth2 introspect URL")
	cmd.Flags().String(authBearerAudienceFlag, "", "OAuth2 audience template")
	cmd.Flags().Bool(versionFlag, false, "Show version")
	cmd.Flags().Bool(debugFlag, false, "debug mode")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.BindPFlags(cmd.Flags())
	viper.AutomaticEnv()
}

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

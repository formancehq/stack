package cmd

import (
	"net/http"

	"github.com/formancehq/ledger/internal/api"
	"github.com/formancehq/ledger/internal/engine"
	"github.com/formancehq/ledger/internal/storage/driver"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/ballast"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	app "github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

const (
	ballastSizeInBytesFlag     = "ballast-size"
	numscriptCacheMaxCountFlag = "numscript-cache-max-count"
	ledgerBatchSizeFlag        = "ledger-batch-size"
	readOnlyFlag               = "read-only"
	autoUpgradeFlag            = "auto-upgrade"
	emitLogsFlag               = "emit-logs"

	ServiceName = "ledger"
)

func NewServe() *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := bindFlagsToViper(cmd); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.New(cmd.OutOrStdout(),
				fx.NopLogger,
				publish.CLIPublisherModule(ServiceName),
				otlptraces.CLITracesModule(),
				otlpmetrics.CLIMetricsModule(),
				auth.CLIAuthModule(),
				driver.CLIModule(cmd),
				engine.Module(engine.Configuration{
					NumscriptCache: engine.NumscriptCacheConfiguration{
						MaxCount: viper.GetInt(numscriptCacheMaxCountFlag),
					},
					GlobalLedgerConfig: engine.GlobalLedgerConfig{
						BatchSize: viper.GetInt(ledgerBatchSizeFlag),
						EmitLogs:  viper.GetBool(emitLogsFlag),
					},
				}),
				ballast.Module(viper.GetUint(ballastSizeInBytesFlag)),
				api.Module(api.Config{
					Version:  Version,
					ReadOnly: viper.GetBool(readOnlyFlag),
				}),
				fx.Invoke(func(lc fx.Lifecycle, driver *driver.Driver) {
					if viper.GetBool(autoUpgradeFlag) {
						lc.Append(fx.Hook{
							OnStart: driver.UpgradeAllBuckets,
						})
					}
				}),
				fx.Invoke(func(lc fx.Lifecycle, h chi.Router, logger logging.Logger) {
					wrappedRouter := chi.NewRouter()
					wrappedRouter.Use(func(handler http.Handler) http.Handler {
						return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
							r = r.WithContext(logging.ContextWithLogger(r.Context(), logger))
							handler.ServeHTTP(w, r)
						})
					})
					wrappedRouter.Use(Log())
					wrappedRouter.Mount("/", h)

					lc.Append(httpserver.NewHook(wrappedRouter, httpserver.WithAddress(viper.GetString(bindFlag))))
				}),
			).Run(cmd.Context())
		},
	}
	cmd.Flags().Uint(ballastSizeInBytesFlag, 0, "Ballast size in bytes, default to 0")
	cmd.Flags().Int(numscriptCacheMaxCountFlag, 1024, "Numscript cache max count")
	cmd.Flags().Int(ledgerBatchSizeFlag, 50, "ledger batch size")
	cmd.Flags().Bool(readOnlyFlag, false, "Read only mode")
	cmd.Flags().Bool(autoUpgradeFlag, false, "Automatically upgrade all schemas")
	cmd.Flags().Bool(emitLogsFlag, false, "Emit logs on NATS")
	return cmd
}

func Log() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(w, r)
			latency := time.Since(start.Time)
			logging.FromContext(r.Context()).WithFields(map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"latency":    latency,
				"user_agent": r.UserAgent(),
				"params":     r.URL.Query().Encode(),
			}).Debug("Request")
		})
	}
}

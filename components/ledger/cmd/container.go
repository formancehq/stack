package cmd

import (
	"crypto/tls"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/cors"
	"github.com/numary/ledger/cmd/internal"
	"github.com/numary/ledger/pkg/api"
	"github.com/numary/ledger/pkg/api/middlewares"
	"github.com/numary/ledger/pkg/api/routes"
	"github.com/numary/ledger/pkg/bus"
	"github.com/numary/ledger/pkg/ledger"
	"github.com/numary/ledger/pkg/redis"
	"github.com/numary/ledger/pkg/storage/sqlstorage"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ServiceName = "ledger"

func resolveOptions(v *viper.Viper, userOptions ...fx.Option) []fx.Option {

	options := make([]fx.Option, 0)

	debug := v.GetBool(service.DebugFlag)
	if debug {
		sqlstorage.InstrumentalizeSQLDrivers()
	}

	options = append(options, publish.CLIPublisherModule(v, ServiceName), bus.LedgerMonitorModule())

	// Handle OpenTelemetry
	options = append(options, otlptraces.CLITracesModule(v))

	switch v.GetString(lockStrategyFlag) {
	case "redis":
		var tlsConfig *tls.Config
		if v.GetBool(lockStrategyRedisTLSEnabledFlag) {
			tlsConfig = &tls.Config{}
			if v.GetBool(lockStrategyRedisTLSInsecureFlag) {
				tlsConfig.InsecureSkipVerify = true
			}
		}
		options = append(options, redis.Module(redis.Config{
			Url:          v.GetString(lockStrategyRedisUrlFlag),
			LockDuration: v.GetDuration(lockStrategyRedisDurationFlag),
			LockRetry:    v.GetDuration(lockStrategyRedisRetryFlag),
			TLSConfig:    tlsConfig,
		}))
	}

	// Handle api part
	options = append(options, api.Module(api.Config{
		StorageDriver: v.GetString(storageDriverFlag),
		Version:       Version,
	}))

	// Handle storage driver
	options = append(options, sqlstorage.DriverModule(sqlstorage.ModuleConfig{
		StorageDriver: v.GetString(storageDriverFlag),
		SQLiteConfig: func() *sqlstorage.SQLiteConfig {
			if v.GetString(storageDriverFlag) != sqlstorage.SQLite.String() {
				return nil
			}
			return &sqlstorage.SQLiteConfig{
				Dir:    v.GetString(storageDirFlag),
				DBName: v.GetString(storageSQLiteDBNameFlag),
			}
		}(),
		PostgresConfig: func() *sqlstorage.PostgresConfig {
			if v.GetString(storageDriverFlag) != sqlstorage.PostgreSQL.String() {
				return nil
			}
			return &sqlstorage.PostgresConfig{
				ConnString: v.GetString(storagePostgresConnectionStringFlag),
			}
		}(),
	}))

	options = append(options, internal.NewAnalyticsModule(v, Version))

	options = append(options, fx.Provide(
		fx.Annotate(func() []ledger.LedgerOption {
			ledgerOptions := []ledger.LedgerOption{}

			if v.GetString(commitPolicyFlag) == "allow-past-timestamps" {
				ledgerOptions = append(ledgerOptions, ledger.WithPastTimestamps)
			}

			return ledgerOptions
		}, fx.ResultTags(ledger.ResolverLedgerOptionsKey)),
	))

	// Handle resolver
	options = append(options, ledger.ResolveModule(
		v.GetInt64(cacheCapacityBytes), v.GetInt64(cacheMaxNumKeys)))

	// Api middlewares
	options = append(options, routes.ProvidePerLedgerMiddleware(func(tp trace.TracerProvider) []gin.HandlerFunc {
		res := make([]gin.HandlerFunc, 0)

		methods := make([]auth.Method, 0)
		if httpBasicMethod := internal.HTTPBasicAuthMethod(v); httpBasicMethod != nil {
			methods = append(methods, httpBasicMethod)
		}
		if len(methods) > 0 {
			res = append(res, func(c *gin.Context) {
				handled := false
				auth.Middleware(methods...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					handled = true
					// The middleware replace the context of the request to include the agent
					// We have to forward it to gin
					c.Request = r
					c.Next()
				})).ServeHTTP(c.Writer, c.Request)
				if !handled {
					c.Abort()
				}
			})
		}
		return res
	}, fx.ParamTags(`optional:"true"`)))

	options = append(options, routes.ProvideMiddlewares(func(tp trace.TracerProvider, logger logging.Logger) []func(handler http.Handler) http.Handler {
		res := make([]func(handler http.Handler) http.Handler, 0)
		res = append(res, cors.New(cors.Options{
			AllowOriginFunc: func(r *http.Request, origin string) bool {
				return true
			},
			AllowCredentials: true,
		}).Handler)
		res = append(res, func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.ServeHTTP(w, r.WithContext(
					logging.ContextWithLogger(r.Context(), logger),
				))
			})
		})
		res = append(res, middlewares.Log())
		return res
	}, fx.ParamTags(`optional:"true"`)))

	return append(options, userOptions...)
}

func NewContainer(v *viper.Viper, userOptions ...fx.Option) *fx.App {
	return fx.New(resolveOptions(v, userOptions...)...)
}

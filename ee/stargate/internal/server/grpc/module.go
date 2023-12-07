package grpc

import (
	"context"
	"net"
	"net/http"

	"github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/server/grpc/interceptors"
	"github.com/formancehq/stack/components/stargate/internal/server/grpc/metrics"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/gogo/status"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func Module(
	bind string,
	jwksURL string,
	maxRetriesJWKSFetching int,
) fx.Option {
	options := make([]fx.Option, 0)

	options = append(options,
		fx.Provide(metrics.RegisterMetricsRegistry),
		fx.Provide(func(logger logging.Logger, metricsRegistry metrics.MetricsRegistry) *interceptors.AuthInterceptor {
			return interceptors.NewAuthInterceptor(logger, jwksURL, maxRetriesJWKSFetching, metricsRegistry)
		}),
		fx.Provide(NewServer),
		fx.Provide(newGrpcServer),
		fx.Invoke(func(lc fx.Lifecycle, grpcServer *grpc.Server, l logging.Logger) error {
			l.Infof("gRPC server listening on %s", bind)
			listener, err := net.Listen("tcp", bind)
			if err != nil {
				return err
			}
			srv := &http2.Server{}

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						for {
							conn, err := listener.Accept()
							if err != nil {
								panic(err)
							}

							go srv.ServeConn(conn, &http2.ServeConnOpts{
								Handler: http.HandlerFunc(grpcServer.ServeHTTP),
							})
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return listener.Close()
				},
			})

			return nil
		}),
	)

	return fx.Options(options...)
}

func newGrpcServer(
	srv *Server,
	l logging.Logger,
	authInterceptor *interceptors.AuthInterceptor,
) *grpc.Server {
	grpcPanicRecoveryHandler := func(p any) (err error) {
		l.Error("panic in grpc server", p)
		return status.Errorf(codes.Internal, "%s", p)
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			grpclogging.StreamServerInterceptor(interceptorLogger(l)),
			authInterceptor.StreamServerInterceptor(),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
	)

	api.RegisterStargateServiceServer(grpcSrv, srv)
	grpc_health_v1.RegisterHealthServer(grpcSrv, health.NewServer())
	reflection.Register(grpcSrv)

	return grpcSrv
}

func interceptorLogger(l logging.Logger) grpclogging.Logger {
	return grpclogging.LoggerFunc(func(ctx context.Context, lvl grpclogging.Level, msg string, fields ...any) {
		switch lvl {
		case grpclogging.LevelDebug:
			var fields []any = append([]any{msg}, fields...)
			l.Debug(fields...)
		case grpclogging.LevelInfo:
			var fields []any = append([]any{msg}, fields...)
			l.Info(fields...)
		case grpclogging.LevelWarn:
			var fields []any = append([]any{msg}, fields...)
			l.Error(fields...)
		case grpclogging.LevelError:
			var fields []any = append([]any{msg}, fields...)
			l.Error(fields...)
		default:
			var fields []any = append([]any{msg}, fields...)
			l.Info(fields...)
		}
	})
}

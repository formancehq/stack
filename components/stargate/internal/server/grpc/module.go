package grpc

import (
	"context"
	"net"
	"time"

	"github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/server/grpc/interceptors"
	"github.com/formancehq/stack/components/stargate/internal/server/grpc/opentelemetry"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/gogo/status"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func Module(
	bind string,
	jwksURL string,
	maxRetriesJWKSFetching int,
) fx.Option {
	options := make([]fx.Option, 0)

	options = append(options,
		fx.Provide(opentelemetry.RegisterMetricsRegistry),
		fx.Provide(func(logger logging.Logger, metricsRegistry opentelemetry.MetricsRegistry) *interceptors.AuthInterceptor {
			return interceptors.NewAuthInterceptor(logger, jwksURL, maxRetriesJWKSFetching, metricsRegistry)
		}),
		fx.Provide(NewServer),
		fx.Provide(newGrpcServer),
		fx.Invoke(func(lc fx.Lifecycle, srv *grpc.Server, l logging.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					listener, err := net.Listen("tcp", bind)
					if err != nil {
						return err
					}

					go func() {
						l.Infof("gRPC server listening on %s", bind)
						err := srv.Serve(listener)
						if err != nil && err != grpc.ErrServerStopped {
							panic(err)
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					srv.GracefulStop()
					srv.Stop()

					return nil
				},
			})
		}),
	)

	return fx.Options(options...)
}

func NewKeepAlivePolicy(
	keepAlivePolicyMinTimeFlag time.Duration,
	keepAlivePolicyPermitWithoutStreamFlag bool,
) keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             keepAlivePolicyMinTimeFlag,
		PermitWithoutStream: keepAlivePolicyPermitWithoutStreamFlag,
	}
}

func NewKeepAliveServerParams(
	keepAliveParamMaxConnectionIdleFlag time.Duration,
	keepAliveParamMaxConnectionAgeFlag time.Duration,
	keepAliveParamMaxConnectionAgeGraceFlag time.Duration,
	keepAliveParamTimeFlag time.Duration,
	keepAliveParamTimeoutFlag time.Duration,
) keepalive.ServerParameters {
	return keepalive.ServerParameters{
		MaxConnectionIdle:     keepAliveParamMaxConnectionIdleFlag,
		MaxConnectionAge:      keepAliveParamMaxConnectionAgeFlag,
		MaxConnectionAgeGrace: keepAliveParamMaxConnectionAgeGraceFlag,
		Time:                  keepAliveParamTimeFlag,
		Timeout:               keepAliveParamTimeoutFlag,
	}
}

func newGrpcServer(
	srv *Server,
	kaep keepalive.EnforcementPolicy,
	kasp keepalive.ServerParameters,
	l logging.Logger,
	authInterceptor *interceptors.AuthInterceptor,
) *grpc.Server {
	grpcPanicRecoveryHandler := func(p any) (err error) {
		l.Error("panic in grpc server", p)
		return status.Errorf(codes.Internal, "%s", p)
	}

	grpcSrv := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
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

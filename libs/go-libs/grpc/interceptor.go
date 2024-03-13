package grpc

// I want to create a logging middle for grpc

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/logging"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func LoggerInterceptor(l logging.Logger) grpclogging.Logger {
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

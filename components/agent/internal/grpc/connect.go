package grpc

import (
	"context"

	membershipgrpc "github.com/formancehq/stack/components/agent/internal/grpc/generated"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func Connect(ctx context.Context, addr string, opts ...grpc.DialOption) (membershipgrpc.ServerClient, error) {
	sharedlogging.FromContext(ctx).Infof("Connect to GRPC server at %s", addr)
	opts = append(opts,
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}
	sharedlogging.FromContext(ctx).Info("Connected to GRPC server!")
	return membershipgrpc.NewServerClient(conn), nil
}

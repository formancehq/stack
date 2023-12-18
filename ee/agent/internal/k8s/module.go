package k8s

import (
	"github.com/formancehq/stack/components/agent/internal/grpc"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Options(
		fx.Provide(fx.Annotate(newClient, fx.As(new(grpc.K8SClient)))),
	)
}

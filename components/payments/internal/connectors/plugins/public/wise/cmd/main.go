package main

import (
	"github.com/formancehq/payments/internal/connectors/grpc"
	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/connectors/plugins/public/wise"
	"github.com/hashicorp/go-plugin"
)

func main() {
	// TODO(polo): metrics
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: grpc.Handshake,
		Plugins: map[string]plugin.Plugin{
			"psp": &grpc.PSPGRPCPlugin{Impl: plugins.NewGRPCImplem(&wise.Plugin{})},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

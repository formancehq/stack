package controllers

import (
	"time"

	"github.com/formancehq/stack/components/stargate/internal/server/http/opentelemetry"
	"github.com/nats-io/nats.go"
)

type StargateControllerConfig struct {
	version string

	natsRequestTimeout time.Duration
}

func NewStargateControllerConfig(
	version string,
	natsRequestTimeout time.Duration,
) StargateControllerConfig {
	return StargateControllerConfig{
		version:            version,
		natsRequestTimeout: natsRequestTimeout,
	}
}

type StargateController struct {
	config StargateControllerConfig

	natsConn        *nats.Conn
	metricsRegistry opentelemetry.MetricsRegistry
}

func NewStargateController(
	natsConn *nats.Conn,
	metricsRegistry opentelemetry.MetricsRegistry,
	config StargateControllerConfig,
) *StargateController {
	return &StargateController{
		natsConn:        natsConn,
		metricsRegistry: metricsRegistry,
		config:          config,
	}
}

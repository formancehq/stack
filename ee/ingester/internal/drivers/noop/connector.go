package noop

import (
	"context"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

type Connector struct{}

func (connector *Connector) Stop(_ context.Context) error {
	return nil
}

func (connector *Connector) Start(_ context.Context) error {
	return nil
}

func (connector *Connector) ClearData(_ context.Context, _ string) error {
	return nil
}

func (connector *Connector) Accept(_ context.Context, logs ...ingester.LogWithModule) ([]error, error) {
	return make([]error, len(logs)), nil
}

func NewConnector(_ drivers.ServiceConfig, _ struct{}, _ logging.Logger) (*Connector, error) {
	return &Connector{}, nil
}

var _ drivers.Driver = (*Connector)(nil)

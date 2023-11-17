package integration_test

import (
	"context"

	"github.com/formancehq/payments/cmd/connectors/internal/integration"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
)

type ConnectorBuilder struct {
	name            string
	uninstall       func(ctx context.Context) error
	resolve         func(descriptor models.TaskDescriptor) task.Task
	install         func(ctx task.ConnectorContext) error
	initiatePayment func(ctx task.ConnectorContext, transfer *models.TransferInitiation) error
}

func (b *ConnectorBuilder) WithUninstall(
	uninstallFunction func(ctx context.Context) error,
) *ConnectorBuilder {
	b.uninstall = uninstallFunction

	return b
}

func (b *ConnectorBuilder) WithResolve(resolveFunction func(name models.TaskDescriptor) task.Task) *ConnectorBuilder {
	b.resolve = resolveFunction

	return b
}

func (b *ConnectorBuilder) WithInstall(installFunction func(ctx task.ConnectorContext) error) *ConnectorBuilder {
	b.install = installFunction

	return b
}

func (b *ConnectorBuilder) Build() integration.Connector {
	return &BuiltConnector{
		name:            b.name,
		uninstall:       b.uninstall,
		resolve:         b.resolve,
		install:         b.install,
		initiatePayment: b.initiatePayment,
	}
}

func NewConnectorBuilder() *ConnectorBuilder {
	return &ConnectorBuilder{}
}

type BuiltConnector struct {
	name            string
	uninstall       func(ctx context.Context) error
	resolve         func(name models.TaskDescriptor) task.Task
	install         func(ctx task.ConnectorContext) error
	initiatePayment func(ctx task.ConnectorContext, transfer *models.TransferInitiation) error
}

func (b *BuiltConnector) SupportedCurrenciesAndDecimals() map[string]int {
	return map[string]int{}
}

func (b *BuiltConnector) InitiatePayment(ctx task.ConnectorContext, transfer *models.TransferInitiation) error {
	if b.initiatePayment != nil {
		return b.initiatePayment(ctx, transfer)
	}

	return nil
}

func (b *BuiltConnector) Name() string {
	return b.name
}

func (b *BuiltConnector) Install(ctx task.ConnectorContext) error {
	if b.install != nil {
		return b.install(ctx)
	}

	return nil
}

func (b *BuiltConnector) Uninstall(ctx context.Context) error {
	if b.uninstall != nil {
		return b.uninstall(ctx)
	}

	return nil
}

func (b *BuiltConnector) Resolve(name models.TaskDescriptor) task.Task {
	if b.resolve != nil {
		return b.resolve(name)
	}

	return nil
}

var _ integration.Connector = &BuiltConnector{}

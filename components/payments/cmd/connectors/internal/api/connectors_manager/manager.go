package connectors_manager

import (
	"context"
	"fmt"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type ConnectorManager struct {
	connector connectors.Connector
	scheduler *task.DefaultTaskScheduler
}

type ConnectorsManager[Config models.ConnectorConfigObject] struct {
	provider         models.ConnectorProvider
	loader           Loader[Config]
	store            Store
	schedulerFactory TaskSchedulerFactory
	publisher        message.Publisher
	messages         *messages.Messages

	connectors map[string]*ConnectorManager
	mu         sync.RWMutex
}

func (l *ConnectorsManager[ConnectorConfig]) logger(ctx context.Context) logging.Logger {
	return logging.FromContext(ctx).WithFields(map[string]interface{}{
		"component": "connector-manager",
		"provider":  l.loader.Name(),
	})
}

func (l *ConnectorsManager[ConnectorConfig]) getManager(connectorID models.ConnectorID) (*ConnectorManager, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	connector, ok := l.connectors[connectorID.String()]
	if !ok {
		return nil, ErrNotInstalled
	}

	return connector, nil
}

func (l *ConnectorsManager[ConnectorConfig]) Connectors() map[string]*ConnectorManager {
	l.mu.RLock()
	defer l.mu.RUnlock()

	copy := make(map[string]*ConnectorManager, len(l.connectors))
	for k, v := range l.connectors {
		copy[k] = v
	}

	return copy
}

func (l *ConnectorsManager[ConnectorConfig]) ReadConfig(
	ctx context.Context,
	connectorID models.ConnectorID,
) (ConnectorConfig, error) {
	var config ConnectorConfig
	connector, err := l.store.GetConnector(ctx, connectorID)
	if err != nil {
		return config, newStorageError(err, "getting connector")
	}

	return l.readConfig(ctx, connector)
}

func (l *ConnectorsManager[ConnectorConfig]) readConfig(
	ctx context.Context,
	connector *models.Connector,
) (ConnectorConfig, error) {
	var config ConnectorConfig
	if connector == nil {
		var err error
		connector, err = l.store.GetConnector(ctx, connector.ID)
		if err != nil {
			return config, newStorageError(err, "getting connector")
		}
	}

	err := connector.ParseConfig(&config)
	if err != nil {
		return config, err
	}

	config = l.loader.ApplyDefaults(config)

	return config, nil
}

func (l *ConnectorsManager[ConnectorConfig]) UpdateConfig(
	ctx context.Context,
	connectorID models.ConnectorID,
	config ConnectorConfig,
) error {
	l.logger(ctx).Infof("Updating config of connector: %s", connectorID)

	connectorManager, err := l.getManager(connectorID)
	if err != nil {
		l.logger(ctx).Errorf("Connector not installed")
		return err
	}

	config = l.loader.ApplyDefaults(config)
	if err = config.Validate(); err != nil {
		return err
	}

	cfg, err := config.Marshal()
	if err != nil {
		return err
	}

	if err := l.store.UpdateConfig(ctx, connectorID, cfg); err != nil {
		return newStorageError(err, "updating connector config")
	}

	// Detach the context since we're launching an async task and we're mostly
	// coming from a HTTP request.
	detachedCtx, span := detachedCtxWithSpan(ctx, trace.SpanFromContext(ctx), "connectorManager.UpdateConfig", connectorID)
	defer span.End()
	if err := connectorManager.connector.UpdateConfig(task.NewConnectorContext(logging.ContextWithLogger(
		detachedCtx,
		logging.FromContext(ctx),
	), connectorManager.scheduler), config); err != nil {
		switch {
		case errors.Is(err, connectors.ErrInvalidConfig):
			return errors.Wrap(ErrValidation, err.Error())
		default:
			return err
		}
	}

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) load(
	ctx context.Context,
	connectorID models.ConnectorID,
	connectorConfig ConnectorConfig,
) error {
	c := l.loader.Load(l.logger(ctx), connectorConfig)
	scheduler := l.schedulerFactory.Make(connectorID, c, l.loader.AllowTasks())

	l.mu.Lock()
	l.connectors[connectorID.String()] = &ConnectorManager{
		connector: c,
		scheduler: scheduler,
	}
	l.mu.Unlock()

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) Install(
	ctx context.Context,
	name string,
	config ConnectorConfig,
) (models.ConnectorID, error) {
	l.logger(ctx).WithFields(map[string]interface{}{
		"config": config,
	}).Infof("Install connector %s", name)

	isInstalled, err := l.store.IsInstalledByConnectorName(ctx, name)
	if err != nil {
		return models.ConnectorID{}, newStorageError(err, "checking if connector is installed")
	}

	if isInstalled {
		l.logger(ctx).Errorf("Connector already installed")
		return models.ConnectorID{}, ErrAlreadyInstalled
	}

	config = l.loader.ApplyDefaults(config)

	if err = config.Validate(); err != nil {
		return models.ConnectorID{}, err
	}

	cfg, err := config.Marshal()
	if err != nil {
		return models.ConnectorID{}, err
	}

	connector := &models.Connector{
		ID: models.ConnectorID{
			Provider:  l.provider,
			Reference: uuid.New(),
		},
		Name:     name,
		Provider: l.provider,
	}

	err = l.store.Install(ctx, connector, cfg)
	if err != nil {
		return models.ConnectorID{}, newStorageError(err, "installing connector")
	}

	if err := l.load(ctx, connector.ID, config); err != nil {
		return models.ConnectorID{}, err
	}

	connectorManager, err := l.getManager(connector.ID)
	if err != nil {
		return models.ConnectorID{}, err
	}

	// Detach the context since we're launching an async task and we're mostly
	// coming from a HTTP request.
	detachedCtx, span := detachedCtxWithSpan(ctx, trace.SpanFromContext(ctx), "connectorManager.Install", connector.ID)
	defer span.End()
	err = connectorManager.connector.Install(task.NewConnectorContext(logging.ContextWithLogger(
		detachedCtx,
		logging.FromContext(ctx),
	), connectorManager.scheduler))
	if err != nil {
		l.logger(ctx).Errorf("Error starting connector: %s", err)

		return models.ConnectorID{}, err
	}

	l.logger(ctx).Infof("Connector installed")

	return connector.ID, nil
}

func (l *ConnectorsManager[ConnectorConfig]) Uninstall(ctx context.Context, connectorID models.ConnectorID) error {
	l.logger(ctx).Infof("Uninstalling connector: %s", connectorID)

	connectorManager, err := l.getManager(connectorID)
	if err != nil {
		l.logger(ctx).Errorf("Connector not installed")
		return err
	}

	err = connectorManager.scheduler.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = connectorManager.connector.Uninstall(ctx)
	if err != nil {
		return err
	}

	err = l.store.Uninstall(ctx, connectorID)
	if err != nil {
		return newStorageError(err, "uninstalling connector")
	}

	if l.publisher != nil {
		err = l.publisher.Publish(events.TopicPayments,
			publish.NewMessage(ctx, l.messages.NewEventResetConnector(connectorID)))
		if err != nil {
			l.logger(ctx).Errorf("Publishing message: %w", err)
		}
	}

	l.mu.Lock()
	delete(l.connectors, connectorID.String())
	l.mu.Unlock()

	l.logger(ctx).Infof("Connector %s uninstalled", connectorID)

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) Restore(ctx context.Context) error {
	l.logger(ctx).Info("Restoring state for all connectors")

	connectors, err := l.store.ListConnectors(ctx)
	if err != nil {
		return newStorageError(err, "listing connectors")
	}

	for _, connector := range connectors {
		if connector.Provider != l.provider {
			continue
		}

		if err := l.restore(ctx, connector); err != nil {
			l.logger(ctx).Errorf("Unable to restore connector %s: %s", connector.Name, err)
			return err
		}
	}

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) restore(ctx context.Context, connector *models.Connector) error {
	l.logger(ctx).Infof("Restoring state for connector: %s", connector.Name)

	if manager, _ := l.getManager(connector.ID); manager != nil {
		return ErrAlreadyRunning
	}

	connectorConfig, err := l.readConfig(ctx, connector)
	if err != nil {
		return err
	}

	if err := l.load(ctx, connector.ID, connectorConfig); err != nil {
		return err
	}

	manager, err := l.getManager(connector.ID)
	if err != nil {
		return err
	}

	if err := manager.scheduler.Restore(ctx); err != nil {
		return err
	}

	l.logger(ctx).Infof("State restored for connector: %s", connector.Name)

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) FindAll(ctx context.Context) ([]*models.Connector, error) {
	connectors, err := l.store.ListConnectors(ctx)
	if err != nil {
		return nil, newStorageError(err, "listing connectors")
	}

	providerConnectors := make([]*models.Connector, 0, len(connectors))
	for _, connector := range connectors {
		if connector.Provider == l.provider {
			providerConnectors = append(providerConnectors, connector)
		}
	}

	return providerConnectors, nil
}

func (l *ConnectorsManager[ConnectorConfig]) IsInstalled(ctx context.Context, connectorID models.ConnectorID) (bool, error) {
	isInstalled, err := l.store.IsInstalledByConnectorID(ctx, connectorID)
	return isInstalled, newStorageError(err, "checking if connector is installed")
}

func (l *ConnectorsManager[ConnectorConfig]) ListTasksStates(
	ctx context.Context,
	connectorID models.ConnectorID,
	q storage.ListTasksQuery,
) (*api.Cursor[models.Task], error) {
	connectorManager, err := l.getManager(connectorID)
	if err != nil {
		return nil, ErrConnectorNotFound
	}

	return connectorManager.scheduler.ListTasks(ctx, q)
}

func (l *ConnectorsManager[Config]) ReadTaskState(ctx context.Context, connectorID models.ConnectorID, taskID uuid.UUID) (*models.Task, error) {
	connectorManager, err := l.getManager(connectorID)
	if err != nil {
		return nil, ErrConnectorNotFound
	}

	return connectorManager.scheduler.ReadTask(ctx, taskID)
}

func (l *ConnectorsManager[ConnectorConfig]) Reset(ctx context.Context, connectorID models.ConnectorID) error {
	connector, err := l.store.GetConnector(ctx, connectorID)
	if err != nil {
		return newStorageError(err, "getting connector")
	}

	config, err := l.readConfig(ctx, connector)
	if err != nil {
		return err
	}

	err = l.Uninstall(ctx, connectorID)
	if err != nil {
		return err
	}

	_, err = l.Install(ctx, connector.Name, config)
	if err != nil {
		return err
	}

	err = l.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, l.messages.NewEventResetConnector(connectorID)))
	if err != nil {
		l.logger(ctx).Errorf("Publishing message: %w", err)
	}

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) InitiatePayment(ctx context.Context, transfer *models.TransferInitiation) error {
	connectorManager, err := l.getManager(transfer.ConnectorID)
	if err != nil {
		return ErrConnectorNotFound
	}

	if err := l.validateAssets(ctx, connectorManager, transfer.ConnectorID, transfer.Asset); err != nil {
		return err
	}

	detachedCtx, span := detachedCtxWithSpan(ctx, trace.SpanFromContext(ctx), "connectorManager.InitiatePayment", transfer.ConnectorID)
	defer span.End()
	err = connectorManager.connector.InitiatePayment(task.NewConnectorContext(detachedCtx, connectorManager.scheduler), transfer)
	if err != nil {
		return fmt.Errorf("initiating transfer: %w", err)
	}

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) ReversePayment(ctx context.Context, transferReversal *models.TransferReversal) error {
	connectorManager, err := l.getManager(transferReversal.ConnectorID)
	if err != nil {
		return ErrConnectorNotFound
	}

	if err := l.validateAssets(ctx, connectorManager, transferReversal.ConnectorID, transferReversal.Asset); err != nil {
		return err
	}

	err = connectorManager.connector.ReversePayment(task.NewConnectorContext(ctx, connectorManager.scheduler), transferReversal)
	if err != nil {
		return fmt.Errorf("reversing transfer: %w", err)
	}

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) CreateExternalBankAccount(ctx context.Context, connectorID models.ConnectorID, bankAccount *models.BankAccount) error {
	connectorManager, err := l.getManager(connectorID)
	if err != nil {
		return ErrConnectorNotFound
	}

	detachedCtx, span := detachedCtxWithSpan(ctx, trace.SpanFromContext(ctx), "connectorManager.CreateExternalBankAccount", connectorID)
	defer span.End()
	err = connectorManager.connector.CreateExternalBankAccount(task.NewConnectorContext(detachedCtx, connectorManager.scheduler), bankAccount)
	if err != nil {
		switch {
		case errors.Is(err, connectors.ErrNotImplemented):
			return errors.Wrap(ErrValidation, "bank account creation not implemented for this connector")
		default:
			return fmt.Errorf("creating bank account: %w", err)
		}
	}

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) CreateWebhookAndContext(
	ctx context.Context,
	webhook *models.Webhook,
) (context.Context, error) {
	connectorManager, err := l.getManager(webhook.ConnectorID)
	if err != nil {
		return nil, ErrConnectorNotFound
	}

	if err := l.store.CreateWebhook(ctx, webhook); err != nil {
		return nil, newStorageError(err, "creating webhook")
	}

	connectorContext := task.NewConnectorContext(ctx, connectorManager.scheduler)
	ctx = task.ContextWithConnectorContext(connectors.ContextWithWebhookID(ctx, webhook.ID), connectorContext)

	return ctx, nil
}

func (l *ConnectorsManager[ConnectorConfig]) validateAssets(
	ctx context.Context,
	connectorManager *ConnectorManager,
	connectorID models.ConnectorID,
	asset models.Asset,
) error {
	supportedCurrencies := connectorManager.connector.SupportedCurrenciesAndDecimals()
	currency, precision, err := models.GetCurrencyAndPrecisionFromAsset(asset)
	if err != nil {
		return errors.Wrap(ErrValidation, err.Error())
	}

	supportedPrecision, ok := supportedCurrencies[currency]
	if !ok {
		return errors.Wrap(ErrValidation, fmt.Sprintf("currency %s not supported", currency))
	}

	if precision != int64(supportedPrecision) {
		return errors.Wrap(ErrValidation, fmt.Sprintf("currency %s has precision %d, but %d is required", currency, precision, supportedPrecision))
	}

	return nil
}

func detachedCtxWithSpan(
	ctx context.Context,
	parentSpan trace.Span,
	spanName string,
	connectorID models.ConnectorID,
) (context.Context, trace.Span) {
	detachedCtx, _ := contextutil.Detached(ctx)

	ctx, span := otel.Tracer().Start(
		detachedCtx,
		spanName,
		trace.WithLinks(trace.Link{
			SpanContext: parentSpan.SpanContext(),
		}),
		trace.WithAttributes(
			attribute.String("connectorID", connectorID.String()),
		),
	)

	return ctx, span
}

func (l *ConnectorsManager[ConnectorConfig]) Close(ctx context.Context) error {
	for _, connectorManager := range l.connectors {
		err := connectorManager.scheduler.Shutdown(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewConnectorManager[ConnectorConfig models.ConnectorConfigObject](
	provider models.ConnectorProvider,
	store Store,
	loader Loader[ConnectorConfig],
	schedulerFactory TaskSchedulerFactory,
	publisher message.Publisher,
	messages *messages.Messages,
) *ConnectorsManager[ConnectorConfig] {
	return &ConnectorsManager[ConnectorConfig]{
		provider:         provider,
		connectors:       make(map[string]*ConnectorManager),
		store:            store,
		loader:           loader,
		schedulerFactory: schedulerFactory,
		publisher:        publisher,
		messages:         messages,
	}
}

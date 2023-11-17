package integration

import (
	"context"
	"fmt"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/cmd/connectors/internal/messages"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrAlreadyInstalled  = errors.New("already installed")
	ErrNotInstalled      = errors.New("not installed")
	ErrNotEnabled        = errors.New("not enabled")
	ErrAlreadyRunning    = errors.New("already running")
	ErrConnectorNotFound = errors.New("connector not found")
	ErrValidation        = errors.New("validation error")
)

type connectorManager struct {
	connector Connector
	scheduler *task.DefaultTaskScheduler
}

type ConnectorsManager[Config models.ConnectorConfigObject] struct {
	provider         models.ConnectorProvider
	loader           Loader[Config]
	store            Repository
	schedulerFactory TaskSchedulerFactory
	publisher        message.Publisher

	connectors map[string]*connectorManager
	mu         sync.RWMutex
}

func (l *ConnectorsManager[ConnectorConfig]) logger(ctx context.Context) logging.Logger {
	return logging.FromContext(ctx).WithFields(map[string]interface{}{
		"component": "connector-manager",
		"provider":  l.loader.Name(),
	})
}

func (l *ConnectorsManager[ConnectorConfig]) getManager(connectorID models.ConnectorID) (*connectorManager, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	connector, ok := l.connectors[connectorID.String()]
	if !ok {
		return nil, ErrNotInstalled
	}

	return connector, nil
}

func (l *ConnectorsManager[ConnectorConfig]) Connectors() map[string]*connectorManager {
	l.mu.RLock()
	defer l.mu.RUnlock()

	copy := make(map[string]*connectorManager, len(l.connectors))
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
		return config, err
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
			return config, err
		}
	}

	err := connector.ParseConfig(&config)
	if err != nil {
		return config, err
	}

	config = l.loader.ApplyDefaults(config)

	return config, nil
}

func (l *ConnectorsManager[ConnectorConfig]) load(
	ctx context.Context,
	connectorID models.ConnectorID,
	connectorConfig ConnectorConfig,
) error {
	c := l.loader.Load(l.logger(ctx), connectorConfig)
	scheduler := l.schedulerFactory.Make(connectorID, c, l.loader.AllowTasks())

	l.mu.Lock()
	l.connectors[connectorID.String()] = &connectorManager{
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
		l.logger(ctx).Errorf("Error checking if connector is installed: %s", err)
		return models.ConnectorID{}, err
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
		return models.ConnectorID{}, err
	}

	if err := l.load(ctx, connector.ID, config); err != nil {
		return models.ConnectorID{}, err
	}

	connectorManager, err := l.getManager(connector.ID)
	if err != nil {
		return models.ConnectorID{}, err
	}

	err = connectorManager.connector.Install(task.NewConnectorContext(logging.ContextWithLogger(
		context.TODO(),
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
		return err
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
		return err
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
		return nil, err
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
	return l.store.IsInstalledByConnectorID(ctx, connectorID)
}

func (l *ConnectorsManager[ConnectorConfig]) ListTasksStates(
	ctx context.Context,
	connectorID models.ConnectorID,
	pagination storage.PaginatorQuery,
) ([]models.Task, storage.PaginationDetails, error) {
	connectorManager, err := l.getManager(connectorID)
	if err != nil {
		return nil, storage.PaginationDetails{}, ErrConnectorNotFound
	}

	return connectorManager.scheduler.ListTasks(ctx, pagination)
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
		return err
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
		publish.NewMessage(ctx, messages.NewEventResetConnector(connectorID)))
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

	err = connectorManager.connector.InitiatePayment(task.NewConnectorContext(ctx, connectorManager.scheduler), transfer)
	if err != nil {
		return fmt.Errorf("initiating transfer: %w", err)
	}

	return nil
}

func (l *ConnectorsManager[ConnectorConfig]) validateAssets(
	ctx context.Context,
	connectorManager *connectorManager,
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
	store Repository,
	loader Loader[ConnectorConfig],
	schedulerFactory TaskSchedulerFactory,
	publisher message.Publisher,
) *ConnectorsManager[ConnectorConfig] {
	return &ConnectorsManager[ConnectorConfig]{
		provider:         provider,
		connectors:       make(map[string]*connectorManager),
		store:            store,
		loader:           loader,
		schedulerFactory: schedulerFactory,
		publisher:        publisher,
	}
}

package integration

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/payments/internal/app/messages"
	"github.com/formancehq/payments/pkg/events"

	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/payments/internal/app/storage"

	"github.com/google/uuid"

	"github.com/formancehq/payments/internal/app/models"

	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyInstalled = errors.New("already installed")
	ErrNotInstalled     = errors.New("not installed")
	ErrNotEnabled       = errors.New("not enabled")
	ErrAlreadyRunning   = errors.New("already running")
)

type ConnectorManager[Config models.ConnectorConfigObject] struct {
	loader           Loader[Config]
	connector        Connector
	store            Repository
	schedulerFactory TaskSchedulerFactory
	scheduler        *task.DefaultTaskScheduler
	publisher        message.Publisher
}

func (l *ConnectorManager[ConnectorConfig]) logger(ctx context.Context) logging.Logger {
	return logging.FromContext(ctx).WithFields(map[string]interface{}{
		"component": "connector-manager",
		"provider":  l.loader.Name(),
	})
}

func (l *ConnectorManager[ConnectorConfig]) Enable(ctx context.Context) error {
	l.logger(ctx).Info("Enabling connector")

	err := l.store.Enable(ctx, l.loader.Name())
	if err != nil {
		return err
	}

	return nil
}

func (l *ConnectorManager[ConnectorConfig]) ReadConfig(ctx context.Context,
) (*ConnectorConfig, error) {
	var config ConnectorConfig

	connector, err := l.store.GetConnector(ctx, l.loader.Name())
	if err != nil {
		return &config, err
	}

	err = connector.ParseConfig(&config)
	if err != nil {
		return &config, err
	}

	config = l.loader.ApplyDefaults(config)

	return &config, nil
}

func (l *ConnectorManager[ConnectorConfig]) load(ctx context.Context, config ConnectorConfig) {
	l.connector = l.loader.Load(l.logger(ctx), config)
	l.scheduler = l.schedulerFactory.Make(l.connector, l.loader.AllowTasks())
}

func (l *ConnectorManager[ConnectorConfig]) Install(ctx context.Context, config ConnectorConfig) error {
	l.logger(ctx).WithFields(map[string]interface{}{
		"config": config,
	}).Infof("Install connector %s", l.loader.Name())

	isInstalled, err := l.store.IsInstalled(ctx, l.loader.Name())
	if err != nil {
		l.logger(ctx).Errorf("Error checking if connector is installed: %s", err)

		return err
	}

	if isInstalled {
		l.logger(ctx).Errorf("Connector already installed")

		return ErrAlreadyInstalled
	}

	config = l.loader.ApplyDefaults(config)

	if err = config.Validate(); err != nil {
		return err
	}

	l.load(ctx, config)

	cfg, err := config.Marshal()
	if err != nil {
		return err
	}

	err = l.store.Install(ctx, l.loader.Name(), cfg)
	if err != nil {
		return err
	}

	err = l.connector.Install(task.NewConnectorContext(logging.ContextWithLogger(
		context.TODO(),
		logging.FromContext(ctx),
	), l.scheduler))
	if err != nil {
		l.logger(ctx).Errorf("Error starting connector: %s", err)

		return err
	}

	l.logger(ctx).Infof("Connector installed")

	return nil
}

func (l *ConnectorManager[ConnectorConfig]) Uninstall(ctx context.Context) error {
	l.logger(ctx).Infof("Uninstalling connector")

	isInstalled, err := l.IsInstalled(ctx)
	if err != nil {
		l.logger(ctx).Errorf("Error checking if connector is installed: %s", err)

		return err
	}

	if !isInstalled {
		l.logger(ctx).Errorf("Connector not installed")

		return ErrNotInstalled
	}

	err = l.scheduler.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = l.connector.Uninstall(ctx)
	if err != nil {
		return err
	}

	err = l.store.Uninstall(ctx, l.loader.Name())
	if err != nil {
		return err
	}

	l.logger(ctx).Info("Connector uninstalled")

	return nil
}

func (l *ConnectorManager[ConnectorConfig]) Restore(ctx context.Context) error {
	l.logger(ctx).Info("Restoring state")

	installed, err := l.IsInstalled(ctx)
	if err != nil {
		return err
	}

	if !installed {
		l.logger(ctx).Info("Not installed, skip")

		return ErrNotInstalled
	}

	enabled, err := l.IsEnabled(ctx)
	if err != nil {
		return err
	}

	if !enabled {
		l.logger(ctx).Info("Not enabled, skip")

		return ErrNotEnabled
	}

	if l.connector != nil {
		return ErrAlreadyRunning
	}

	config, err := l.ReadConfig(ctx)
	if err != nil {
		return err
	}

	l.load(ctx, *config)

	err = l.scheduler.Restore(ctx)
	if err != nil {
		l.logger(ctx).Errorf("Unable to restore scheduler: %s", err)

		return err
	}

	l.logger(ctx).Info("State restored")

	return nil
}

func (l *ConnectorManager[ConnectorConfig]) Disable(ctx context.Context) error {
	l.logger(ctx).Info("Disabling connector")

	return l.store.Disable(ctx, l.loader.Name())
}

func (l *ConnectorManager[ConnectorConfig]) IsEnabled(ctx context.Context) (bool, error) {
	return l.store.IsEnabled(ctx, l.loader.Name())
}

func (l *ConnectorManager[ConnectorConfig]) FindAll(ctx context.Context) ([]*models.Connector, error) {
	return l.store.ListConnectors(ctx)
}

func (l *ConnectorManager[ConnectorConfig]) IsInstalled(ctx context.Context) (bool, error) {
	return l.store.IsInstalled(ctx, l.loader.Name())
}

func (l *ConnectorManager[ConnectorConfig]) ListTasksStates(ctx context.Context, pagination storage.Paginator,
) ([]models.Task, storage.PaginationDetails, error) {
	return l.scheduler.ListTasks(ctx, pagination)
}

func (l *ConnectorManager[Config]) ReadTaskState(ctx context.Context, taskID uuid.UUID) (*models.Task, error) {
	return l.scheduler.ReadTask(ctx, taskID)
}

func (l *ConnectorManager[ConnectorConfig]) Reset(ctx context.Context) error {
	config, err := l.ReadConfig(ctx)
	if err != nil {
		return err
	}

	err = l.Uninstall(ctx)
	if err != nil {
		return err
	}

	err = l.Install(ctx, *config)
	if err != nil {
		return err
	}

	err = l.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, messages.NewEventResetConnector(l.loader.Name())))
	if err != nil {
		l.logger(ctx).Errorf("Publishing message: %w", err)
	}

	return nil
}

type Transfer struct {
	Source      string
	Destination string
	Currency    string
	Amount      int64
}

func (l *ConnectorManager[ConnectorConfig]) InitiateTransfer(ctx context.Context, transfer Transfer) (uuid.UUID, error) {
	newTransfer, err := l.store.CreateNewTransfer(ctx, l.loader.Name(),
		transfer.Source, transfer.Destination, transfer.Currency, transfer.Amount)
	if err != nil {
		return uuid.Nil, fmt.Errorf("creating new transfer: %w", err)
	}

	err = l.connector.InitiateTransfer(task.NewConnectorContext(ctx, l.scheduler), newTransfer)
	if err != nil {
		return uuid.Nil, fmt.Errorf("initiating transfer: %w", err)
	}

	return newTransfer.ID, nil
}

func (l *ConnectorManager[ConnectorConfig]) ListTransfers(ctx context.Context) ([]models.Transfer, error) {
	transfers, err := l.store.ListTransfers(ctx, l.loader.Name())
	if err != nil {
		return nil, fmt.Errorf("retrieving transfers: %w", err)
	}

	return transfers, nil
}

func NewConnectorManager[ConnectorConfig models.ConnectorConfigObject](
	store Repository,
	loader Loader[ConnectorConfig],
	schedulerFactory TaskSchedulerFactory,
	publisher message.Publisher,
) *ConnectorManager[ConnectorConfig] {
	return &ConnectorManager[ConnectorConfig]{
		store:            store,
		loader:           loader,
		schedulerFactory: schedulerFactory,
		publisher:        publisher,
	}
}

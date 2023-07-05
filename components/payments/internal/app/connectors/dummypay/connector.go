package dummypay

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/app/integration"

	"github.com/formancehq/payments/internal/app/models"

	"github.com/formancehq/payments/internal/app/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

// Name is the name of the connector.
const Name = models.ConnectorProviderDummyPay

// Connector is the connector for the dummy payment connector.
type Connector struct {
	logger logging.Logger
	cfg    Config
	fs     fs
}

func (c *Connector) InitiateTransfer(ctx task.ConnectorContext, transfer models.Transfer) error {
	// TODO implement me
	panic("implement me")
}

// Install executes post-installation steps to read and generate files.
// It is called after the connector is installed.
func (c *Connector) Install(ctx task.ConnectorContext) error {
	readFilesDescriptor, err := models.EncodeTaskDescriptor(newTaskReadFiles())
	if err != nil {
		return fmt.Errorf("failed to create read files task descriptor: %w", err)
	}

	if err = ctx.Scheduler().Schedule(ctx.Context(), readFilesDescriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		// No need to restart this task, since the connector is not existing or
		// was uninstalled previously, the task does not exists in the database
		Restart: false,
	}); err != nil {
		return fmt.Errorf("failed to schedule task to read files: %w", err)
	}

	generateFilesDescriptor, err := models.EncodeTaskDescriptor(newTaskGenerateFiles())
	if err != nil {
		return fmt.Errorf("failed to create generate files task descriptor: %w", err)
	}

	if err = ctx.Scheduler().Schedule(ctx.Context(), generateFilesDescriptor, models.TaskSchedulerOptions{
		ScheduleOption: models.OPTIONS_RUN_NOW,
		Restart:        false,
	}); err != nil {
		return fmt.Errorf("failed to schedule task to generate files: %w", err)
	}

	return nil
}

// Uninstall executes pre-uninstallation steps to remove the generated files.
// It is called before the connector is uninstalled.
func (c *Connector) Uninstall(ctx context.Context) error {
	c.logger.Infof("Removing generated files from '%s'...", c.cfg.Directory)

	err := removeFiles(c.cfg, c.fs)
	if err != nil {
		return fmt.Errorf("failed to remove generated files: %w", err)
	}

	return nil
}

// Resolve resolves a task execution request based on the task descriptor.
func (c *Connector) Resolve(descriptor models.TaskDescriptor) task.Task {
	taskDescriptor, err := models.DecodeTaskDescriptor[TaskDescriptor](descriptor)
	if err != nil {
		panic(err)
	}

	c.logger.Infof("Executing '%s' task...", taskDescriptor.Key)

	return handleResolve(c.cfg, taskDescriptor, c.fs)
}

var _ integration.Connector = &Connector{}

func newConnector(logger logging.Logger, cfg Config, fs fs) *Connector {
	return &Connector{
		logger: logger.WithFields(map[string]any{
			"component": "connector",
		}),
		cfg: cfg,
		fs:  fs,
	}
}

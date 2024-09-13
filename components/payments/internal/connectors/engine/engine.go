package engine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/connectors/engine/plugins"
	"github.com/formancehq/payments/internal/connectors/engine/webhooks"
	"github.com/formancehq/payments/internal/connectors/engine/workflow"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type Engine interface {
	InstallConnector(ctx context.Context, provider string, rawConfig json.RawMessage) (models.ConnectorID, error)
	UninstallConnector(ctx context.Context, connectorID models.ConnectorID) error
	CreateBankAccount(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error)
	HandleWebhook(ctx context.Context, urlPath string, webhook models.Webhook) error

	OnStart(ctx context.Context) error
}

type engine struct {
	temporalClient client.Client

	workers  *Workers
	plugins  plugins.Plugins
	storage  storage.Storage
	webhooks webhooks.Webhooks
}

func New(temporalClient client.Client, workers *Workers, plugins plugins.Plugins, storage storage.Storage) Engine {
	return &engine{
		temporalClient: temporalClient,
		workers:        workers,
		plugins:        plugins,
		storage:        storage,
	}
}

func (e *engine) InstallConnector(ctx context.Context, provider string, rawConfig json.RawMessage) (models.ConnectorID, error) {
	config := models.DefaultConfig()
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		return models.ConnectorID{}, err
	}

	if err := config.Validate(); err != nil {
		return models.ConnectorID{}, errors.Wrap(ErrValidation, err.Error())
	}

	connectorID := models.ConnectorID{
		Reference: config.Name,
		Provider:  provider,
	}

	err := e.plugins.RegisterPlugin(connectorID)
	if err != nil {
		return models.ConnectorID{}, handlePluginError(err)
	}

	err = e.workers.AddWorker(connectorID)
	if err != nil {
		return models.ConnectorID{}, err
	}

	// Launch the workflow
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("install-%s", connectorID.String()),
			TaskQueue:                                connectorID.Reference,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
		},
		workflow.RunInstallConnector,
		workflow.InstallConnector{
			ConnectorID: connectorID,
			RawConfig:   rawConfig,
		},
	)
	if err != nil {
		return models.ConnectorID{}, err
	}

	// Wait for installation to complete
	if err := run.Get(ctx, nil); err != nil {
		return models.ConnectorID{}, err
	}

	return connectorID, nil
}

func (e *engine) UninstallConnector(ctx context.Context, connectorID models.ConnectorID) error {
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("uninstall-%s", connectorID.String()),
			TaskQueue:                                connectorID.Reference,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
		},
		workflow.RunUninstallConnector,
		workflow.UninstallConnector{
			ConnectorID: connectorID,
		},
	)
	if err != nil {
		return err
	}

	// Wait for uninstallation to complete
	if err := run.Get(ctx, nil); err != nil {
		return err
	}

	if err := e.workers.RemoveWorker(connectorID); err != nil {
		return err
	}

	if err := e.plugins.UnregisterPlugin(connectorID); err != nil {
		return handlePluginError(err)
	}

	return nil
}

func (e *engine) CreateBankAccount(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error) {
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("create-bank-account-%s-%s", connectorID.String(), bankAccountID.String()),
			TaskQueue:                                connectorID.Reference,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
		},
		workflow.RunCreateBankAccount,
		workflow.CreateBankAccount{
			ConnectorID:   connectorID,
			BankAccountID: bankAccountID,
		},
	)
	if err != nil {
		return nil, err
	}

	var bankAccount models.BankAccount
	// Wait for bank account creation to complete
	if err := run.Get(ctx, &bankAccount); err != nil {
		return nil, err
	}

	return &bankAccount, nil
}

func (e *engine) HandleWebhook(ctx context.Context, urlPath string, webhook models.Webhook) error {
	return e.webhooks.HandleWebhook(ctx, urlPath, webhook)
}

func (e *engine) OnStart(ctx context.Context) error {
	query := storage.NewListConnectorsQuery(
		bunpaginate.NewPaginatedQueryOptions(storage.ConnectorQuery{}).
			WithPageSize(100),
	)

	for {
		connectors, err := e.storage.ConnectorsList(ctx, query)
		if err != nil {
			return err
		}

		for _, connector := range connectors.Data {
			if err := e.onStartPlugin(ctx, connector); err != nil {
				return err
			}
		}

		if !connectors.HasMore {
			break
		}

		err = bunpaginate.UnmarshalCursor(connectors.Next, &query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *engine) onStartPlugin(ctx context.Context, connector models.Connector) error {
	err := e.plugins.RegisterPlugin(connector.ID)
	if err != nil {
		return err
	}

	err = e.workers.AddWorker(connector.ID)
	if err != nil {
		return err
	}

	// Launch the workflow
	_, err = e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("install-%s", connector.ID.String()),
			TaskQueue:                                connector.ID.Reference,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
		},
		workflow.RunInstallConnector,
		workflow.InstallConnector{
			ConnectorID: connector.ID,
			RawConfig:   connector.Config,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

var _ Engine = &engine{}

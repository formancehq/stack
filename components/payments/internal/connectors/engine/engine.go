package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/contextutil"
	"github.com/formancehq/payments/internal/connectors/engine/plugins"
	"github.com/formancehq/payments/internal/connectors/engine/webhooks"
	"github.com/formancehq/payments/internal/connectors/engine/workflow"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

type Engine interface {
	InstallConnector(ctx context.Context, provider string, rawConfig json.RawMessage) (models.ConnectorID, error)
	UninstallConnector(ctx context.Context, connectorID models.ConnectorID) error
	ResetConnector(ctx context.Context, connectorID models.ConnectorID) error
	ForwardBankAccount(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error)
	HandleWebhook(ctx context.Context, urlPath string, webhook models.Webhook) error
	CreatePool(ctx context.Context, pool models.Pool) error
	AddAccountToPool(ctx context.Context, id uuid.UUID, accountID models.AccountID) error
	RemoveAccountFromPool(ctx context.Context, id uuid.UUID, accountID models.AccountID) error
	DeletePool(ctx context.Context, poolID uuid.UUID) error

	OnStart(ctx context.Context) error
}

type engine struct {
	temporalClient client.Client

	workers  *Workers
	plugins  plugins.Plugins
	storage  storage.Storage
	webhooks webhooks.Webhooks

	stack string
}

func New(temporalClient client.Client, workers *Workers, plugins plugins.Plugins, storage storage.Storage, webhooks webhooks.Webhooks, stack string) Engine {
	return &engine{
		temporalClient: temporalClient,
		workers:        workers,
		plugins:        plugins,
		storage:        storage,
		webhooks:       webhooks,
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

	connector := models.Connector{
		ID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  provider,
		},
		Name:      config.Name,
		CreatedAt: time.Now().UTC(),
		Provider:  provider,
		Config:    rawConfig,
	}

	if err := e.storage.ConnectorsInstall(ctx, connector); err != nil {
		return models.ConnectorID{}, err
	}

	err := e.plugins.RegisterPlugin(connector.ID)
	if err != nil {
		return models.ConnectorID{}, handlePluginError(err)
	}

	err = e.workers.AddWorker(connector.ID.String())
	if err != nil {
		return models.ConnectorID{}, err
	}

	// Launch the workflow
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("install-%s", connector.ID.String()),
			TaskQueue:                                connector.ID.String(),
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
		},
		workflow.RunInstallConnector,
		workflow.InstallConnector{
			ConnectorID: connector.ID,
			RawConfig:   rawConfig,
			Config:      config,
		},
	)
	if err != nil {
		return models.ConnectorID{}, err
	}

	// Wait for installation to complete
	if err := run.Get(ctx, nil); err != nil {
		return models.ConnectorID{}, err
	}

	return connector.ID, nil
}

func (e *engine) UninstallConnector(ctx context.Context, connectorID models.ConnectorID) error {
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("uninstall-%s", connectorID.String()),
			TaskQueue:                                connectorID.String(),
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
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

	if err := e.workers.RemoveWorker(connectorID.String()); err != nil {
		return err
	}

	if err := e.plugins.UnregisterPlugin(connectorID); err != nil {
		return handlePluginError(err)
	}

	return nil
}

func (e *engine) ResetConnector(ctx context.Context, connectorID models.ConnectorID) error {
	connector, err := e.storage.ConnectorsGet(ctx, connectorID)
	if err != nil {
		return err
	}

	// Detached the context to avoid being in a weird state if request is
	// cancelled in the middle of the operation.
	ctx, _ = contextutil.Detached(ctx)
	if err := e.UninstallConnector(ctx, connectorID); err != nil {
		return err
	}

	_, err = e.InstallConnector(ctx, connectorID.Provider, connector.Config)
	if err != nil {
		return err
	}

	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                    fmt.Sprintf("reset-%s", connectorID.String()),
			TaskQueue:             connectorID.String(),
			WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
		},
		workflow.RunSendEvents,
		workflow.SendEvents{
			ConnectorReset: &connectorID,
		},
	)
	if err != nil {
		return err
	}

	if err := run.Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func (e *engine) ForwardBankAccount(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error) {
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("create-bank-account-%s-%s", connectorID.String(), bankAccountID.String()),
			TaskQueue:                                connectorID.String(),
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
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
	configs, err := e.webhooks.GetConfigs(webhook.ConnectorID, urlPath)
	if err != nil {
		return err
	}

	var config *models.WebhookConfig
	for _, c := range configs {
		if !strings.Contains(urlPath, c.URLPath) {
			continue
		}

		config = &c
		break
	}

	if config == nil {
		return errors.New("webhook config not found")
	}

	ctx, _ = contextutil.Detached(ctx)
	if _, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("webhook-%s-%s", webhook.ConnectorID.String(), webhook.ID),
			TaskQueue:                                webhook.ConnectorID.String(),
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
		},
		workflow.RunHandleWebhooks,
		workflow.HandleWebhooks{
			ConnectorID:   webhook.ConnectorID,
			WebhookConfig: *config,
			Webhook:       webhook,
		},
	); err != nil {
		return err
	}

	return nil
}

func (e *engine) CreatePool(ctx context.Context, pool models.Pool) error {
	if err := e.storage.PoolsUpsert(ctx, pool); err != nil {
		return err
	}

	ctx, _ = contextutil.Detached(ctx)
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("pools-creation-%s", pool.IdempotencyKey()),
			TaskQueue:                                defaultWorkerName,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
		},
		workflow.RunSendEvents,
		workflow.SendEvents{
			PoolsCreation: &pool,
		},
	)
	if err != nil {
		return err
	}

	if err := run.Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func (e *engine) AddAccountToPool(ctx context.Context, id uuid.UUID, accountID models.AccountID) error {
	if err := e.storage.PoolsAddAccount(ctx, id, accountID); err != nil {
		return err
	}

	ctx, _ = contextutil.Detached(ctx)
	pool, err := e.storage.PoolsGet(ctx, id)
	if err != nil {
		return err
	}

	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("pools-add-account-%s", pool.IdempotencyKey()),
			TaskQueue:                                defaultWorkerName,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
		},
		workflow.RunSendEvents,
		workflow.SendEvents{
			PoolsCreation: pool,
		},
	)
	if err != nil {
		return err
	}

	if err := run.Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func (e *engine) RemoveAccountFromPool(ctx context.Context, id uuid.UUID, accountID models.AccountID) error {
	if err := e.storage.PoolsRemoveAccount(ctx, id, accountID); err != nil {
		return err
	}

	ctx, _ = contextutil.Detached(ctx)
	pool, err := e.storage.PoolsGet(ctx, id)
	if err != nil {
		return err
	}

	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("pools-remove-account-%s", pool.IdempotencyKey()),
			TaskQueue:                                defaultWorkerName,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
		},
		workflow.RunSendEvents,
		workflow.SendEvents{
			PoolsCreation: pool,
		},
	)
	if err != nil {
		return err
	}

	if err := run.Get(ctx, nil); err != nil {
		return err
	}

	return nil
}

func (e *engine) DeletePool(ctx context.Context, poolID uuid.UUID) error {
	if err := e.storage.PoolsDelete(ctx, poolID); err != nil {
		return err
	}
	ctx, _ = contextutil.Detached(ctx)
	run, err := e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("pools-deletion-%s", poolID.String()),
			TaskQueue:                                defaultWorkerName,
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
		},
		workflow.RunSendEvents,
		workflow.SendEvents{
			PoolsDeletion: &poolID,
		},
	)
	if err != nil {
		return err
	}

	if err := run.Get(ctx, nil); err != nil {
		return err
	}

	return nil
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

	err = e.workers.AddWorker(connector.ID.String())
	if err != nil {
		return err
	}

	config := models.DefaultConfig()
	if err := json.Unmarshal(connector.Config, &config); err != nil {
		return err
	}

	// Launch the workflow
	_, err = e.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:                                       fmt.Sprintf("install-%s", connector.ID.String()),
			TaskQueue:                                connector.ID.String(),
			WorkflowIDReusePolicy:                    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
			WorkflowExecutionErrorWhenAlreadyStarted: false,
			SearchAttributes: map[string]interface{}{
				workflow.SearchAttributeStack: e.stack,
			},
		},
		workflow.RunInstallConnector,
		workflow.InstallConnector{
			ConnectorID: connector.ID,
			RawConfig:   connector.Config,
			Config:      config,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

var _ Engine = &engine{}

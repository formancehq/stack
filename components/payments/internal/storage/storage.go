package storage

import (
	"context"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Storage interface {
	// Accounts
	AccountsUpsert(ctx context.Context, accounts []models.Account) error
	AccountsGet(ctx context.Context, id models.AccountID) (*models.Account, error)
	AccountsList(ctx context.Context, q ListAccountsQuery) (*bunpaginate.Cursor[models.Account], error)
	AccountsDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Balances
	BalancesUpsert(ctx context.Context, balances []models.Balance) error
	BalancesDeleteForConnectorID(ctx context.Context, connectorID models.ConnectorID) error
	BalancesList(ctx context.Context, q ListBalancesQuery) (*bunpaginate.Cursor[models.Balance], error)
	BalancesGetAt(ctx context.Context, accountID models.AccountID, at time.Time) ([]*models.Balance, error)

	// Bank Accounts
	BankAccountsUpsert(ctx context.Context, bankAccount models.BankAccount) error
	BankAccountsUpdateMetadata(ctx context.Context, id uuid.UUID, metadata map[string]string) error
	BankAccountsGet(ctx context.Context, id uuid.UUID, expand bool) (*models.BankAccount, error)
	BankAccountsList(ctx context.Context, q ListBankAccountsQuery) (*bunpaginate.Cursor[models.BankAccount], error)
	BankAccountsAddRelatedAccount(ctx context.Context, relatedAccount models.BankAccountRelatedAccount) error
	BankAccountsDeleteRelatedAccountFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Connectors
	ConnectorsInstall(ctx context.Context, c models.Connector) error
	ConnectorsUninstall(ctx context.Context, id models.ConnectorID) error
	ConnectorsGet(ctx context.Context, id models.ConnectorID) (*models.Connector, error)
	ConnectorsList(ctx context.Context, q ListConnectorsQuery) (*bunpaginate.Cursor[models.Connector], error)

	// Payments
	PaymentsUpsert(ctx context.Context, payments []models.Payment) error
	PaymentsUpdateMetadata(ctx context.Context, id models.PaymentID, metadata map[string]string) error
	PaymentsGet(ctx context.Context, id models.PaymentID) (*models.Payment, error)
	PaymentsList(ctx context.Context, q ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error)
	PaymentsDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Pools
	PoolsUpsert(ctx context.Context, pool models.Pool) error
	PoolsGet(ctx context.Context, id uuid.UUID) (*models.Pool, error)
	PoolsDelete(ctx context.Context, id uuid.UUID) error
	PoolsAddAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error
	PoolsRemoveAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error
	PoolsList(ctx context.Context, q ListPoolsQuery) (*bunpaginate.Cursor[models.Pool], error)

	// Schedules
	SchedulesUpsert(ctx context.Context, schedule models.Schedule) error
	SchedulesList(ctx context.Context, q ListSchedulesQuery) (*bunpaginate.Cursor[models.Schedule], error)
	SchedulesGet(ctx context.Context, id string, connectorID models.ConnectorID) (*models.Schedule, error)
	SchedulesDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// State
	StatesUpsert(ctx context.Context, state models.State) error
	StatesGet(ctx context.Context, id models.StateID) (models.State, error)
	StatesDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Tasks
	TasksUpsert(ctx context.Context, connectorID models.ConnectorID, tasks models.Tasks) error
	TasksGet(ctx context.Context, connectorID models.ConnectorID) (*models.Tasks, error)
	TasksDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Webhooks Configs
	WebhooksConfigsUpsert(ctx context.Context, webhooksConfigs []models.WebhookConfig) error
	WebhooksConfigsDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Webhooks
	WebhooksInsert(ctx context.Context, webhook models.Webhook) error
	WebhooksDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error

	// Workflow Instances
	InstancesUpsert(ctx context.Context, instance models.Instance) error
	InstancesUpdate(ctx context.Context, instance models.Instance) error
	InstancesList(ctx context.Context, q ListInstancesQuery) (*bunpaginate.Cursor[models.Instance], error)
	InstancesDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error
}

const encryptionOptions = "compress-algo=1, cipher-algo=aes256"

type store struct {
	db                  *bun.DB
	configEncryptionKey string
}

func newStorage(db *bun.DB, configEncryptionKey string) Storage {
	return &store{db: db, configEncryptionKey: configEncryptionKey}
}

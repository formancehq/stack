package backend

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/google/uuid"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package backend . Backend
type Backend interface {
	// Accounts
	AccountsList(ctx context.Context, query storage.ListAccountsQuery) (*bunpaginate.Cursor[models.Account], error)
	AccountsGet(ctx context.Context, id models.AccountID) (*models.Account, error)
	BankAccountsList(ctx context.Context, query storage.ListBankAccountsQuery) (*bunpaginate.Cursor[models.BankAccount], error)
	BankAccountsUpdateMetadata(ctx context.Context, id uuid.UUID, metadata map[string]string) error
	BankAccountsForwardToConnector(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error)

	// Balances
	BalancesList(ctx context.Context, query storage.ListBalancesQuery) (*bunpaginate.Cursor[models.Balance], error)
	PoolsBalancesAt(ctx context.Context, poolID uuid.UUID, at time.Time) ([]models.AggregatedBalance, error)

	// Bank Accounts
	BankAccountsCreate(ctx context.Context, bankAccount models.BankAccount) error
	BankAccountsGet(ctx context.Context, id uuid.UUID) (*models.BankAccount, error)

	// Connectors
	ConnectorsConfigs() plugins.Configs
	ConnectorsConfig(ctx context.Context, connectorID models.ConnectorID) (json.RawMessage, error)
	ConnectorsList(ctx context.Context, query storage.ListConnectorsQuery) (*bunpaginate.Cursor[models.Connector], error)
	ConnectorsInstall(ctx context.Context, provider string, config json.RawMessage) (models.ConnectorID, error)
	ConnectorsUninstall(ctx context.Context, connectorID models.ConnectorID) error
	ConnectorsReset(ctx context.Context, connectorID models.ConnectorID) error

	// Payments
	PaymentsUpdateMetadata(ctx context.Context, id models.PaymentID, metadata map[string]string) error
	PaymentsList(ctx context.Context, query storage.ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error)
	PaymentsGet(ctx context.Context, id models.PaymentID) (*models.Payment, error)

	// Pools
	PoolsCreate(ctx context.Context, pool models.Pool) error
	PoolsGet(ctx context.Context, id uuid.UUID) (*models.Pool, error)
	PoolsList(ctx context.Context, query storage.ListPoolsQuery) (*bunpaginate.Cursor[models.Pool], error)
	PoolsDelete(ctx context.Context, id uuid.UUID) error
	PoolsAddAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error
	PoolsRemoveAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error

	// Schedules
	SchedulesList(ctx context.Context, query storage.ListSchedulesQuery) (*bunpaginate.Cursor[models.Schedule], error)
	SchedulesGet(ctx context.Context, id string, connectorID models.ConnectorID) (*models.Schedule, error)

	// Webhooks
	ConnectorsHandleWebhooks(ctx context.Context, urlPath string, webhook models.Webhook) error

	// Workflows Instances
	WorkflowsInstancesList(ctx context.Context, query storage.ListInstancesQuery) (*bunpaginate.Cursor[models.Instance], error)
}

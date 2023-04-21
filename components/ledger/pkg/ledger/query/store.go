package query

import (
	"context"

	"github.com/formancehq/ledger/pkg/core"
)

type Store interface {
	UpdateNextLogID(ctx context.Context, u uint64) error
	IsInitialized() bool
	GetNextLogID(ctx context.Context) (uint64, error)
	ReadLogsRange(ctx context.Context, idMin, idMax uint64) ([]core.PersistedLog, error)
	GetAccountWithVolumes(ctx context.Context, address string) (*core.AccountWithVolumes, error)
	GetTransaction(ctx context.Context, id uint64) (*core.ExpandedTransaction, error)
	UpdateAccountsMetadata(ctx context.Context, update []core.Account) error
	InsertTransactions(ctx context.Context, insert ...core.ExpandedTransaction) error
	UpdateTransactionsMetadata(ctx context.Context, update ...core.TransactionWithMetadata) error
	EnsureAccountsExist(ctx context.Context, accounts []string) error
	UpdateVolumes(ctx context.Context, moreRecentLogID uint64, update ...core.AccountsAssetsVolumes) error
}

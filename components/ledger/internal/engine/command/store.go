package command

import (
	"context"
	"math/big"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/machine/vm"
)

type Store interface {
	vm.Store
	InsertLogs(ctx context.Context, logs ...*ledger.ChainedLog) error
	GetLastLog(ctx context.Context) (*ledger.ChainedLog, error)
	ReadLogWithIdempotencyKey(ctx context.Context, key string) (*ledger.ChainedLog, error)
	ReadLastLogWithType(ctx context.Context, logType ...ledger.LogType) (*ledger.ChainedLog, error)
	GetTransactionByReference(ctx context.Context, ref string) (*ledger.ExpandedTransaction, error)
	GetTransaction(ctx context.Context, txID *big.Int) (*ledger.Transaction, error)
}

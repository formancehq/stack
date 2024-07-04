package chain

import (
	"context"

	ledger "github.com/formancehq/ledger/internal"
)

type Store interface {
	GetLastLog(ctx context.Context) (*ledger.ChainedLog, error)
	GetLastTransaction(ctx context.Context) (*ledger.ExpandedTransaction, error)
}

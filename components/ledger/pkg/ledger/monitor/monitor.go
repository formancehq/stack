package monitor

import (
	"context"
	"fmt"

	"github.com/formancehq/ledger/pkg/core"
)

type Monitor interface {
	CommittedTransactions(ctx context.Context, ledger string, res ...core.ExpandedTransaction)
	SavedMetadata(ctx context.Context, ledger, targetType, id string, metadata core.Metadata)
	RevertedTransaction(ctx context.Context, ledger string, reverted, revert *core.ExpandedTransaction)
}

type noOpMonitor struct{}

func (n noOpMonitor) CommittedTransactions(ctx context.Context, s string, res ...core.ExpandedTransaction) {
	fmt.Println("FAKE MONITOR COMMITED TRANSACTIONS")
}
func (n noOpMonitor) SavedMetadata(ctx context.Context, ledger string, targetType string, id string, metadata core.Metadata) {
	fmt.Println("FAKE MONITOR SAVE METADATA")
}
func (n noOpMonitor) RevertedTransaction(ctx context.Context, ledger string, reverted, revert *core.ExpandedTransaction) {
	fmt.Println("FAKE MONITOR REVETED TRANSACTION")
}

var _ Monitor = &noOpMonitor{}

func NewNoOpMonitor() *noOpMonitor {
	return &noOpMonitor{}
}

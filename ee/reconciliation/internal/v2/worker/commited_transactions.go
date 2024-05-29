package worker

import (
	"context"

	"github.com/formancehq/reconciliation/internal/v2/models"
)

func (l *Listener) handleCommittedTransactions(
	ctx context.Context,
	policy models.Policy,
	rule *RuleTree,
	ev transaction,
) error {

	return nil
}

package worker

import (
	"context"

	"github.com/formancehq/reconciliation/internal/v2/models"
)

func (l *Listener) handleCommittedTransactions(
	ctx context.Context,
	policy models.Policy,
	rules []*models.Rule,
	ev transaction,
) error {
	return nil
}

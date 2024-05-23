package send

import (
	"github.com/formancehq/orchestration/internal/workflow/activities"
	"github.com/formancehq/orchestration/internal/workflow/stages/internal"
	"go.temporal.io/sdk/workflow"
)

func RunUpdate(ctx workflow.Context, update Update) (err error) {
	switch {
	case update.Account != nil:
		return activities.AddAccountMetadata(internal.InfiniteRetryContext(ctx), activities.AddAccountMetadataRequest{
			Ledger:   update.Account.Ledger,
			Account:  update.Account.ID,
			Metadata: update.Account.Metadata,
		})
	default:
		panic("invalid update specification")
	}
}

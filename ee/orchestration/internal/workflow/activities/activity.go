package activities

import (
	"context"
	"fmt"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	client *sdk.Formance
}

func (a Activities) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Name: "GetAccount",
			Func: a.GetAccount,
		}).
		Append(temporalworker.Definition{
			Name: "AddAccountMetadata",
			Func: a.AddAccountMetadata,
		}).
		Append(temporalworker.Definition{
			Name: "CreateTransaction",
			Func: a.CreateTransaction,
		}).
		Append(temporalworker.Definition{
			Name: "StripeTransfer",
			Func: a.StripeTransfer,
		}).
		Append(temporalworker.Definition{
			Name: "GetPayment",
			Func: a.GetPayment,
		}).
		Append(temporalworker.Definition{
			Name: "ConfirmHold",
			Func: a.ConfirmHold,
		}).
		Append(temporalworker.Definition{
			Name: "CreditWallet",
			Func: a.CreditWallet,
		}).
		Append(temporalworker.Definition{
			Name: "DebitWallet",
			Func: a.DebitWallet,
		}).
		Append(temporalworker.Definition{
			Name: "GetWallet",
			Func: a.GetWallet,
		}).
		Append(temporalworker.Definition{
			Name: "ListWallets",
			Func: a.ListWallets,
		}).
		Append(temporalworker.Definition{
			Name: "VoidHold",
			Func: a.VoidHold,
		})
}

func New(client *sdk.Formance) Activities {
	return Activities{
		client: client,
	}
}

func executeActivity(ctx workflow.Context, activity any, ret any, request any) error {
	if err := workflow.ExecuteActivity(ctx, activity, request).Get(ctx, ret); err != nil {
		var timeoutError *temporal.TimeoutError
		if errors.As(err, &timeoutError) {
			return errors.New(timeoutError.Message())
		}
		return err
	}
	return nil
}

func getLedgerIK(ctx context.Context) *string {
	activityInfo := activity.GetInfo(ctx)
	return pointer.For(fmt.Sprintf("%s-%s", activityInfo.WorkflowExecution.RunID, activityInfo.ActivityID))
}

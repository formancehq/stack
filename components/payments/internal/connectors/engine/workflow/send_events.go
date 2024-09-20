package workflow

import (
	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

type SendEvents struct {
	Account        *models.Account
	Balance        *models.Balance
	BankAccount    *models.BankAccount
	Payment        *models.Payment
	ConnectorReset *models.ConnectorID
	PoolsCreation  *models.Pool
	PoolsDeletion  *uuid.UUID
}

func (w Workflow) runSendEvents(
	ctx workflow.Context,
	sendEvents SendEvents,
) error {
	if sendEvents.Account != nil {
		err := activities.EventsSendAccount(
			infiniteRetryContext(ctx),
			*sendEvents.Account,
		)
		if err != nil {
			return err
		}
	}

	if sendEvents.Balance != nil {
		err := activities.EventsSendBalance(
			infiniteRetryContext(ctx),
			*sendEvents.Balance,
		)
		if err != nil {
			return err
		}
	}

	if sendEvents.BankAccount != nil {
		err := activities.EventsSendBankAccount(
			infiniteRetryContext(ctx),
			*sendEvents.BankAccount,
		)
		if err != nil {
			return err
		}
	}

	if sendEvents.Payment != nil {
		for _, adjustment := range sendEvents.Payment.Adjustments {
			err := activities.EventsSendPayment(
				infiniteRetryContext(ctx),
				*sendEvents.Payment,
				adjustment,
			)
			if err != nil {
				return err
			}
		}
	}

	if sendEvents.ConnectorReset != nil {
		err := activities.EventsSendConnectorReset(
			infiniteRetryContext(ctx),
			*sendEvents.ConnectorReset,
		)
		if err != nil {
			return err
		}
	}

	if sendEvents.PoolsCreation != nil {
		err := activities.EventsSendPoolCreation(
			infiniteRetryContext(ctx),
			*sendEvents.PoolsCreation,
		)
		if err != nil {
			return err
		}
	}

	if sendEvents.PoolsDeletion != nil {
		err := activities.EventsSendPoolDeletion(
			infiniteRetryContext(ctx),
			*sendEvents.PoolsDeletion,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

var RunSendEvents any

func init() {
	RunSendEvents = Workflow{}.runSendEvents
}

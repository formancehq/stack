package ingestion

import (
	"context"

	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

func (i *DefaultIngester) LinkBankAccountWithAccount(ctx context.Context, bankAccount *models.BankAccount, accountID *models.AccountID) error {
	if err := i.store.LinkBankAccountWithAccount(ctx, bankAccount.ID, accountID); err != nil {
		return err
	}

	bankAccount.AccountID = accountID

	if err := i.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			messages.NewEventSavedBankAccounts(bankAccount),
		),
	); err != nil {
		return err
	}

	return nil
}

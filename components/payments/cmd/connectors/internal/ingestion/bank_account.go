package ingestion

import (
	"context"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/google/uuid"
)

func (i *DefaultIngester) LinkBankAccountWithAccount(ctx context.Context, bankAccount *models.BankAccount, accountID *models.AccountID) error {
	adjustment := &models.BankAccountAdjustment{
		ID:            uuid.New(),
		CreatedAt:     time.Now().UTC(),
		BankAccountID: bankAccount.ID,
		ConnectorID:   accountID.ConnectorID,
		AccountID:     *accountID,
	}

	if err := i.store.AddBankAccountAdjustment(ctx, adjustment); err != nil {
		return err
	}

	bankAccount.Adjustments = append(bankAccount.Adjustments, adjustment)

	if err := i.publisher.Publish(
		events.TopicPayments,
		publish.NewMessage(
			ctx,
			i.messages.NewEventSavedBankAccounts(bankAccount),
		),
	); err != nil {
		return err
	}

	return nil
}

package messages

import (
	"time"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/stack/libs/events"
	"github.com/formancehq/stack/libs/events/payments"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewEventSavedAccounts(accounts []models.Account) (*events.Event, error) {
	accs := make([]*payments.Account, len(accounts))
	for i, account := range accounts {
		accs[i] = &payments.Account{
			Id:        account.ID.String(),
			CreatedAt: timestamppb.New(account.CreatedAt),
			Reference: account.Reference,
			Provider:  account.Provider,
			Type:      account.Type.String(),
		}
	}

	now := time.Now().UTC()

	return &events.Event{
		CreatedAt: timestamppb.New(now),
		App:       EventApp,
		Event: &events.Event_AccountSaved{
			AccountSaved: &payments.AccountSaved{
				Accounts: accs,
			},
		},
	}, nil
}

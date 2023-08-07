package messages

import (
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/pkg/events"
)

type balanceMessagePayload struct {
	CreatedAt time.Time `json:"createdAt"`
	AccountID string    `json:"accountID"`
	Asset     string    `json:"asset"`
	Balance   *big.Int  `json:"balance"`
}

func NewEventSavedBalances(balances []*models.Balance) events.EventMessage {
	payload := make([]balanceMessagePayload, len(balances))

	for balanceIdx, balance := range balances {
		payload[balanceIdx] = balanceMessagePayload{
			CreatedAt: balance.CreatedAt,
			AccountID: balance.AccountID.String(),
			Asset:     balance.Asset.String(),
			Balance:   balance.Balance,
		}
	}

	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedBalances,
		Payload: payload,
	}
}

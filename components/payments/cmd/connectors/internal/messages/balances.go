package messages

import (
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type balanceMessagePayload struct {
	CreatedAt time.Time `json:"createdAt"`
	Provider  string    `json:"provider"`
	AccountID string    `json:"accountID"`
	Asset     string    `json:"asset"`
	Balance   *big.Int  `json:"balance"`
}

func NewEventSavedBalances(balance *models.Balance, provider models.ConnectorProvider) events.EventMessage {
	payload := balanceMessagePayload{
		CreatedAt: balance.CreatedAt,
		Provider:  provider.String(),
		AccountID: balance.AccountID.String(),
		Asset:     balance.Asset.String(),
		Balance:   balance.Balance,
	}

	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedBalances,
		Payload: payload,
	}
}

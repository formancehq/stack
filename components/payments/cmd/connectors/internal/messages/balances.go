package messages

import (
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type balanceMessagePayload struct {
	AccountID   string    `json:"accountID"`
	ConnectorID string `json:"connectorId"`
	CreatedAt   time.Time `json:"createdAt"`
	Asset       string    `json:"asset"`
	Balance     *big.Int  `json:"balance"`
}

func NewEventSavedBalances(balance *models.Balance) events.EventMessage {
	payload := balanceMessagePayload{
		CreatedAt:   balance.CreatedAt,
		ConnectorID: balance.ConnectorID.String(),
		AccountID:   balance.AccountID.String(),
		Asset:       balance.Asset.String(),
		Balance:     balance.Balance,
	}

	return events.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedBalances,
		Payload: payload,
	}
}

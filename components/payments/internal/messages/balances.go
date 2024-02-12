package messages

import (
	"math/big"
	"time"

	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type balanceMessagePayload struct {
	AccountID   string    `json:"accountID"`
	ConnectorID string    `json:"connectorId"`
	Provider    string    `json:"provider"`
	CreatedAt   time.Time `json:"createdAt"`
	Asset       string    `json:"asset"`
	Balance     *big.Int  `json:"balance"`
}

func (m *Messages) NewEventSavedBalances(balance *models.Balance) publish.EventMessage {
	payload := balanceMessagePayload{
		CreatedAt:   balance.CreatedAt,
		ConnectorID: balance.ConnectorID.String(),
		Provider:    balance.ConnectorID.Provider.String(),
		AccountID:   balance.AccountID.String(),
		Asset:       balance.Asset.String(),
		Balance:     balance.Balance,
	}

	return publish.EventMessage{
		Date:    time.Now().UTC(),
		App:     events.EventApp,
		Version: events.EventVersion,
		Type:    events.EventTypeSavedBalances,
		Payload: payload,
	}
}

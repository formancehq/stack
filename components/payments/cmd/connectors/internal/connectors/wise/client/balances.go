package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Balance struct {
	ID       uint64 `json:"id"`
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Amount   struct {
		Value    json.Number `json:"value"`
		Currency string      `json:"currency"`
	} `json:"amount"`
	ReservedAmount struct {
		Value    json.Number `json:"value"`
		Currency string      `json:"currency"`
	} `json:"reservedAmount"`
	CashAmount struct {
		Value    json.Number `json:"value"`
		Currency string      `json:"currency"`
	} `json:"cashAmount"`
	TotalWorth struct {
		Value    json.Number `json:"value"`
		Currency string      `json:"currency"`
	} `json:"totalWorth"`
	CreationTime     time.Time `json:"creationTime"`
	ModificationTime time.Time `json:"modificationTime"`
	Visible          bool      `json:"visible"`
}

func (w *Client) GetBalances(ctx context.Context, profileID uint64) ([]*Balance, error) {
	f := connectors.ClientMetrics(ctx, "wise", "list_balances")
	now := time.Now()
	defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, w.endpoint(fmt.Sprintf("v4/profiles/%d/balances?types=STANDARD", profileID)), http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, unmarshalError(res.StatusCode, res.Body).Error()
	}

	var balances []*Balance
	err = json.NewDecoder(res.Body).Decode(&balances)
	if err != nil {
		return nil, fmt.Errorf("failed to decode account: %w", err)
	}

	return balances, nil
}

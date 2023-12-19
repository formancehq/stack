package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/google/uuid"
)

type Quote struct {
	ID uuid.UUID `json:"id"`
}

func (w *Client) CreateQuote(ctx context.Context, profileID, currency string, amount *big.Float) (Quote, error) {
	f := connectors.ClientMetrics(ctx, "wise", "create_quote")
	now := time.Now()
	defer f(ctx, now)

	var response Quote

	req, err := json.Marshal(map[string]interface{}{
		"sourceCurrency": currency,
		"targetCurrency": currency,
		"sourceAmount":   amount,
	})
	if err != nil {
		return response, err
	}

	res, err := w.httpClient.Post(w.endpoint("v3/profiles/"+profileID+"/quotes"), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return response, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return response, unmarshalError(res.StatusCode, res.Body).Error()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("failed to get response from quote: %w", err)
	}

	return response, nil
}

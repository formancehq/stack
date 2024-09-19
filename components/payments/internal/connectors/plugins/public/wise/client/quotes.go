package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
	"github.com/google/uuid"
)

type Quote struct {
	ID uuid.UUID `json:"id"`
}

func (c *Client) CreateQuote(ctx context.Context, profileID, currency string, amount json.Number) (Quote, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "wise", "create_quote")
	// now := time.Now()
	// defer f(ctx, now)

	var quote Quote

	reqBody, err := json.Marshal(map[string]interface{}{
		"sourceCurrency": currency,
		"targetCurrency": currency,
		"sourceAmount":   amount,
	})
	if err != nil {
		return quote, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint("v3/profiles/"+profileID+"/quotes"),
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return quote, err
	}
	req.Header.Set("Content-Type", "application/json")

	var errRes wiseErrors
	statusCode, err := c.httpClient.Do(req, &quote, &errRes)
	switch err {
	case nil:
		return quote, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return quote, errRes.Error(statusCode).Error()
	}
	return quote, fmt.Errorf("failed to get response from quote: %w", err)
}

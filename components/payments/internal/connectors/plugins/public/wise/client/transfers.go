package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

type Transfer struct {
	ID             uint64      `json:"id"`
	Reference      string      `json:"reference"`
	Status         string      `json:"status"`
	SourceAccount  uint64      `json:"sourceAccount"`
	SourceCurrency string      `json:"sourceCurrency"`
	SourceValue    json.Number `json:"sourceValue"`
	TargetAccount  uint64      `json:"targetAccount"`
	TargetCurrency string      `json:"targetCurrency"`
	TargetValue    json.Number `json:"targetValue"`
	Business       uint64      `json:"business"`
	Created        string      `json:"created"`
	//nolint:tagliatelle // allow for clients
	CustomerTransactionID string `json:"customerTransactionId"`
	Details               struct {
		Reference string `json:"reference"`
	} `json:"details"`
	Rate float64 `json:"rate"`
	User uint64  `json:"user"`

	SourceBalanceID      uint64 `json:"-"`
	DestinationBalanceID uint64 `json:"-"`

	CreatedAt time.Time `json:"-"`
}

func (t *Transfer) UnmarshalJSON(data []byte) error {
	type Alias Transfer

	aux := &struct {
		Created string `json:"created"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error

	t.CreatedAt, err = time.Parse("2006-01-02 15:04:05", aux.Created)
	if err != nil {
		return fmt.Errorf("failed to parse created time: %w", err)
	}

	return nil
}

func (c *Client) GetTransfers(ctx context.Context, profileID uint64, offset int, limit int) ([]Transfer, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "wise", "list_transfers")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, c.endpoint("v1/transfers"), http.NoBody)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("limit", fmt.Sprintf("%d", limit))
	q.Add("profile", fmt.Sprintf("%d", profileID))
	q.Add("offset", fmt.Sprintf("%d", offset))
	req.URL.RawQuery = q.Encode()

	var transfers []Transfer
	var errRes wiseErrors
	statusCode, err := c.httpClient.Do(req, &transfers, &errRes)
	switch err {
	case nil:
		// fallthrough
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return transfers, errRes.Error(statusCode).Error()
	default:
		return transfers, fmt.Errorf("failed to get transfers: %w", err)
	}

	for i, transfer := range transfers {
		var sourceProfileID, targetProfileID uint64
		if transfer.SourceAccount != 0 {
			recipientAccount, err := c.GetRecipientAccount(ctx, transfer.SourceAccount)
			if err != nil {
				return nil, fmt.Errorf("failed to get source profile id: %w", err)
			}

			sourceProfileID = recipientAccount.Profile
		}

		if transfer.TargetAccount != 0 {
			recipientAccount, err := c.GetRecipientAccount(ctx, transfer.TargetAccount)
			if err != nil {
				return nil, fmt.Errorf("failed to get target profile id: %w", err)
			}

			targetProfileID = recipientAccount.Profile
		}

		// TODO(polo): fetching balances for each transfer is not efficient
		// and can be quite long. We should consider caching balances, but
		// at the same time we will develop a feature soon to get balances
		// for every accounts, so caching is not a solution.
		switch {
		case sourceProfileID == 0 && targetProfileID == 0:
			// Do nothing
		case sourceProfileID == targetProfileID && sourceProfileID != 0:
			// Same profile id for target and source
			balances, err := c.GetBalances(ctx, sourceProfileID)
			if err != nil {
				return nil, fmt.Errorf("failed to get balances: %w", err)
			}
			for _, balance := range balances {
				if balance.Currency == transfer.SourceCurrency {
					transfers[i].SourceBalanceID = balance.ID
				}

				if balance.Currency == transfer.TargetCurrency {
					transfers[i].DestinationBalanceID = balance.ID
				}
			}
		default:
			if sourceProfileID != 0 {
				balances, err := c.GetBalances(ctx, sourceProfileID)
				if err != nil {
					return nil, fmt.Errorf("failed to get balances: %w", err)
				}
				for _, balance := range balances {
					if balance.Currency == transfer.SourceCurrency {
						transfers[i].SourceBalanceID = balance.ID
					}
				}
			}

			if targetProfileID != 0 {
				balances, err := c.GetBalances(ctx, targetProfileID)
				if err != nil {
					return nil, fmt.Errorf("failed to get balances: %w", err)
				}
				for _, balance := range balances {
					if balance.Currency == transfer.TargetCurrency {
						transfers[i].DestinationBalanceID = balance.ID
					}
				}
			}

		}
	}
	return transfers, nil
}

func (c *Client) GetTransfer(ctx context.Context, transferID string) (*Transfer, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "wise", "get_transfer")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, c.endpoint("v1/transfers/"+transferID), http.NoBody)
	if err != nil {
		return nil, err
	}

	var transfer Transfer
	var errRes wiseErrors
	statusCode, err := c.httpClient.Do(req, &transfer, &errRes)
	switch err {
	case nil:
		return &transfer, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error(statusCode).Error()
	}
	return nil, fmt.Errorf("failed to get transfer: %w", err)
}

func (c *Client) CreateTransfer(ctx context.Context, quote Quote, targetAccount uint64, transactionID string) (*Transfer, error) {
	// TODO(polo): metrics
	// metrics.GetMetricsRegistry().ConnectorPSPCalls().Add(ctx, 1, metric.WithAttributes([]attribute.KeyValue{
	// 	attribute.String("connector", "wise"),
	// 	attribute.String("operation", "initiate_transfer"),
	// }...))

	reqBody, err := json.Marshal(map[string]interface{}{
		"targetAccount":         targetAccount,
		"quoteUuid":             quote.ID.String(),
		"customerTransactionId": transactionID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost, c.endpoint("v1/transfers"), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	var transfer Transfer
	var errRes wiseErrors
	statusCode, err := c.httpClient.Do(req, &transfer, &errRes)
	switch err {
	case nil:
		return &transfer, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error(statusCode).Error()
	}
	return nil, fmt.Errorf("failed to create transfer: %w", err)
}

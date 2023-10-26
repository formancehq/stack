package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

func (w *Client) GetTransfers(ctx context.Context, profile *Profile) ([]Transfer, error) {
	var transfers []Transfer

	limit := 10
	offset := 0

	for {
		req, err := http.NewRequestWithContext(ctx,
			http.MethodGet, w.endpoint("v1/transfers"), http.NoBody)
		if err != nil {
			return transfers, err
		}

		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", limit))
		q.Add("profile", fmt.Sprintf("%d", profile.ID))
		q.Add("offset", fmt.Sprintf("%d", offset))
		req.URL.RawQuery = q.Encode()

		res, err := w.httpClient.Do(req)
		if err != nil {
			return transfers, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, unmarshalError(res.StatusCode, res.Body).Error()
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var transferList []Transfer

		err = json.Unmarshal(body, &transferList)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal transfers: %w", err)
		}

		for i, transfer := range transferList {
			var sourceProfileID, targetProfileID uint64
			if transfer.SourceAccount != 0 {
				recipientAccount, err := w.GetRecipientAccount(ctx, transfer.SourceAccount)
				if err != nil {
					return nil, fmt.Errorf("failed to get source profile id: %w", err)
				}

				sourceProfileID = recipientAccount.Profile
			}

			if transfer.TargetAccount != 0 {
				recipientAccount, err := w.GetRecipientAccount(ctx, transfer.TargetAccount)
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
				balances, err := w.GetBalances(ctx, sourceProfileID)
				if err != nil {
					return nil, fmt.Errorf("failed to get balances: %w", err)
				}
				for _, balance := range balances {
					if balance.Currency == transfer.SourceCurrency {
						transferList[i].SourceBalanceID = balance.ID
					}

					if balance.Currency == transfer.TargetCurrency {
						transferList[i].DestinationBalanceID = balance.ID
					}
				}
			default:
				if sourceProfileID != 0 {
					balances, err := w.GetBalances(ctx, sourceProfileID)
					if err != nil {
						return nil, fmt.Errorf("failed to get balances: %w", err)
					}
					for _, balance := range balances {
						if balance.Currency == transfer.SourceCurrency {
							transferList[i].SourceBalanceID = balance.ID
						}
					}
				}

				if targetProfileID != 0 {
					balances, err := w.GetBalances(ctx, targetProfileID)
					if err != nil {
						return nil, fmt.Errorf("failed to get balances: %w", err)
					}
					for _, balance := range balances {
						if balance.Currency == transfer.TargetCurrency {
							transferList[i].DestinationBalanceID = balance.ID
						}
					}
				}

			}
		}

		transfers = append(transfers, transferList...)

		if len(transferList) < limit {
			break
		}

		offset += limit
	}

	return transfers, nil
}

func (w *Client) GetTransfer(ctx context.Context, transferID string) (*Transfer, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, w.endpoint("v1/transfers/"+transferID), http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		res.Body.Close()

		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err = res.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	var transfer Transfer
	err = json.Unmarshal(body, &transfer)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfer: %w", err)
	}

	return &transfer, nil
}

func (w *Client) CreateTransfer(quote Quote, targetAccount uint64, transactionID string) (*Transfer, error) {
	req, err := json.Marshal(map[string]interface{}{
		"targetAccount":         targetAccount,
		"quoteUuid":             quote.ID.String(),
		"customerTransactionId": transactionID,
	})
	if err != nil {
		return nil, err
	}

	res, err := w.httpClient.Post(w.endpoint("v1/transfers"), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response Transfer
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from transfer: %w", err)
	}

	return &response, nil
}

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

type Payout struct {
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

func (t *Payout) UnmarshalJSON(data []byte) error {
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

func (w *Client) GetPayout(ctx context.Context, payoutID string) (*Payout, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, w.endpoint("v1/transfers/"+payoutID), http.NoBody)
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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var payout Payout
	err = json.Unmarshal(body, &payout)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfer: %w", err)
	}

	return &payout, nil
}

func (w *Client) CreatePayout(quote Quote, targetAccount uint64, transactionID string) (*Payout, error) {
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
		return nil, unmarshalError(res.StatusCode, res.Body).Error()
	}

	var response Payout
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to get response from transfer: %w", err)
	}

	return &response, nil
}

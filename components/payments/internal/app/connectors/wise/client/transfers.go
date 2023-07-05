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
	ID             uint64  `json:"id"`
	Reference      string  `json:"reference"`
	Status         string  `json:"status"`
	SourceAccount  uint64  `json:"sourceAccount"`
	SourceCurrency string  `json:"sourceCurrency"`
	SourceValue    float64 `json:"sourceValue"`
	TargetAccount  uint64  `json:"targetAccount"`
	TargetCurrency string  `json:"targetCurrency"`
	TargetValue    float64 `json:"targetValue"`
	Business       uint64  `json:"business"`
	Created        string  `json:"created"`
	//nolint:tagliatelle // allow for clients
	CustomerTransactionID string `json:"customerTransactionId"`
	Details               struct {
		Reference string `json:"reference"`
	} `json:"details"`
	Rate float64 `json:"rate"`
	User uint64  `json:"user"`

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

		body, err := io.ReadAll(res.Body)
		if err != nil {
			res.Body.Close()

			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if err = res.Body.Close(); err != nil {
			return nil, fmt.Errorf("failed to close response body: %w", err)
		}

		var transferList []Transfer

		err = json.Unmarshal(body, &transferList)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal transfers: %w", err)
		}

		transfers = append(transfers, transferList...)

		if len(transferList) < limit {
			break
		}

		offset += limit
	}

	return transfers, nil
}

func (w *Client) CreateTransfer(quote Quote, targetAccount uint64, transactionID string) error {
	req, err := json.Marshal(map[string]interface{}{
		"targetAccount":         targetAccount,
		"quoteUuid":             quote.ID.String(),
		"customerTransactionId": transactionID,
	})
	if err != nil {
		return err
	}

	res, err := w.httpClient.Post(w.endpoint("v1/transfers"), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

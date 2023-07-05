package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/uuid"
)

type Quote struct {
	ID uuid.UUID `json:"id"`
}

func (w *Client) CreateQuote(profileID uint64, currency string, amount int64) (Quote, error) {
	var response Quote

	req, err := json.Marshal(map[string]interface{}{
		"sourceCurrency": currency,
		"targetCurrency": currency,
		"sourceAmount":   amount,
	})
	if err != nil {
		return response, err
	}

	res, err := w.httpClient.Post(w.endpoint("v3/profiles/"+fmt.Sprint(profileID)+"/quotes"), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return response, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body: %w", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, fmt.Errorf("failed to unmarshal profiles: %w", err)
	}

	return response, nil
}

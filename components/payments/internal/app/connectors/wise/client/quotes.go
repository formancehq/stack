package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/google/uuid"
)

type Quote struct {
	ID uuid.UUID `json:"id"`
}

func (w *Client) CreateQuote(profileID string, currency string, amount *big.Float) (Quote, error) {
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

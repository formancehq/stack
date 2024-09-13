package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PayinResponse struct {
	ID               string `json:"Id"`
	Tag              string `json:"Tag"`
	CreationDate     int64  `json:"CreationDate"`
	ResultCode       string `json:"ResultCode"`
	ResultMessage    string `json:"ResultMessage"`
	AuthorId         string `json:"AuthorId"`
	CreditedUserId   string `json:"CreditedUserId"`
	DebitedFunds     Funds  `json:"DebitedFunds"`
	CreditedFunds    Funds  `json:"CreditedFunds"`
	Fees             Funds  `json:"Fees"`
	Status           string `json:"Status"`
	ExecutionDate    int64  `json:"ExecutionDate"`
	Type             string `json:"Type"`
	CreditedWalletID string `json:"CreditedWalletId"`
	PaymentType      string `json:"PaymentType"`
	ExecutionType    string `json:"ExecutionType"`
}

func (c *Client) GetPayin(ctx context.Context, payinID string) (*PayinResponse, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "mangopay", "get_payin")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/payins/%s", c.endpoint, c.clientID, payinID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get payin request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get payin: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			// TODO(polo): log error
			_ = err
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var payinResponse PayinResponse
	if err := json.NewDecoder(resp.Body).Decode(&payinResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payin response body: %w", err)
	}

	return &payinResponse, nil
}

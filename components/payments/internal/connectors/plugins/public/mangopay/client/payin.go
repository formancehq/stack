package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
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

	var payinResponse PayinResponse
	_, err = c.httpClient.Do(req, &payinResponse, nil)
	switch err {
	case nil:
		return &payinResponse, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, err
	}
	return nil, fmt.Errorf("failed to get payin response: %w", err)
}

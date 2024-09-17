package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

type Funds struct {
	Currency string      `json:"Currency"`
	Amount   json.Number `json:"Amount"`
}

type TransferRequest struct {
	AuthorID         string `json:"AuthorId"`
	CreditedUserID   string `json:"CreditedUserId,omitempty"`
	DebitedFunds     Funds  `json:"DebitedFunds"`
	Fees             Funds  `json:"Fees"`
	DebitedWalletID  string `json:"DebitedWalletId"`
	CreditedWalletID string `json:"CreditedWalletId"`
}

type TransferResponse struct {
	ID               string `json:"Id"`
	CreationDate     int64  `json:"CreationDate"`
	AuthorID         string `json:"AuthorId"`
	CreditedUserID   string `json:"CreditedUserId"`
	DebitedFunds     Funds  `json:"DebitedFunds"`
	Fees             Funds  `json:"Fees"`
	CreditedFunds    Funds  `json:"CreditedFunds"`
	Status           string `json:"Status"`
	ResultCode       string `json:"ResultCode"`
	ResultMessage    string `json:"ResultMessage"`
	Type             string `json:"Type"`
	ExecutionDate    int64  `json:"ExecutionDate"`
	Nature           string `json:"Nature"`
	DebitedWalletID  string `json:"DebitedWalletId"`
	CreditedWalletID string `json:"CreditedWalletId"`
}

func (c *Client) GetWalletTransfer(ctx context.Context, transferID string) (TransferResponse, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "mangopay", "get_transfer")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/transfers/%s", c.endpoint, c.clientID, transferID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return TransferResponse{}, fmt.Errorf("failed to create login request: %w", err)
	}

	var transfer TransferResponse
	_, err = c.httpClient.Do(req, &transfer, nil)
	switch err {
	case nil:
		return transfer, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return transfer, err
	}
	return transfer, fmt.Errorf("failed to get transfer response: %w", err)
}

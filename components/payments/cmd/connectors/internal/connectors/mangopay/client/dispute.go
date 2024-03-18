package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Dispute struct {
	Id                       string `json:"Id"`
	Tag                      string `json:"Tag"`
	InitialTransactionId     string `json:"InitialTransactionId"`
	InitialTransactionType   string `json:"InitialTransactionType"`
	InitialTransactionNature string `json:"InitialTransactionNature"`
	DisputeType              string `json:"DisputeType"`
	ContestDeadlineDate      int64  `json:"ContestDeadlineDate"`
	DisputedFunds            Funds  `json:"DisputedFunds"`
	ContestedFunds           Funds  `json:"ContestedFunds"`
	Status                   string `json:"Status"`
	StatusMessage            string `json:"StatusMessage"`
	DisputeReason            string `json:"DisputeReason"`
	ResultCode               string `json:"ResultCode"`
	ResultMessage            string `json:"ResultMessage"`
	CreationDate             int64  `json:"CreationDate"`
	ClosedDate               int64  `json:"ClosedDate"`
}

func (c *Client) GetDispute(ctx context.Context, disputeID string) (*Dispute, error) {
	f := connectors.ClientMetrics(ctx, "mangopay", "get_dispute")
	now := time.Now()
	defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/disputes/%s", c.endpoint, c.clientID, disputeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get dispute request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get dispute: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalErrorWithRetry(resp.StatusCode, resp.Body).Error()
	}

	var dispute Dispute
	if err := json.NewDecoder(resp.Body).Decode(&dispute); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispute response body: %w", err)
	}

	return &dispute, nil
}

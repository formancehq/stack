package client

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Beneficiary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created string `json:"created"`
}

func (m *Client) GetBeneficiaries(ctx context.Context) ([]*Beneficiary, error) {
	f := connectors.ClientMetrics(ctx, "modulr", "list_beneficiaries")
	now := time.Now()
	defer f(ctx, now)

	resp, err := m.httpClient.Get(m.buildEndpoint("beneficiaries"))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res responseWrapper[[]*Beneficiary]
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Content, nil
}

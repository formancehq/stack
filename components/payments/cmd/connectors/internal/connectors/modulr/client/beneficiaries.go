package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Beneficiary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created string `json:"created"`
}

func (m *Client) GetBeneficiaries(ctx context.Context, page, pageSize int, modifiedSince string) (*responseWrapper[[]*Beneficiary], error) {
	f := connectors.ClientMetrics(ctx, "modulr", "list_beneficiaries")
	now := time.Now()
	defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.buildEndpoint("beneficiaries"), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(pageSize))
	if modifiedSince != "" {
		q.Add("modifiedSince", modifiedSince)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := m.httpClient.Do(req)
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

	return &res, nil
}

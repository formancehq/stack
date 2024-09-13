package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Beneficiary struct {
	ID                    string    `json:"id"`
	BankAccountHolderName string    `json:"bank_account_holder_name"`
	Name                  string    `json:"name"`
	Currency              string    `json:"currency"`
	CreatedAt             time.Time `json:"created_at"`
	// Contains a lot more fields that will be not used on our side for now
}

func (c *Client) GetBeneficiaries(ctx context.Context, page int, pageSize int) ([]*Beneficiary, int, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "currencycloud", "list_beneficiaries")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.buildEndpoint("v2/beneficiaries/find"), http.NoBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page", fmt.Sprint(page))
	q.Add("per_page", fmt.Sprint(pageSize))
	q.Add("order", "created_at")
	q.Add("order_asc_desc", "asc")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	//nolint:tagliatelle // allow for client code
	type response struct {
		Beneficiaries []*Beneficiary `json:"beneficiaries"`
		Pagination    struct {
			NextPage int `json:"next_page"`
		} `json:"pagination"`
	}

	var res response
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, 0, err
	}

	return res.Beneficiaries, res.Pagination.NextPage, nil
}

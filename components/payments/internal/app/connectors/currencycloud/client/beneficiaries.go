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

func (c *Client) GetBeneficiaries(ctx context.Context, page int) ([]*Beneficiary, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.buildEndpoint("v2/beneficiaries/find?page=%d&per_page=25", page), http.NoBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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

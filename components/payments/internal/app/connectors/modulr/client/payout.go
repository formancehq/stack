package client

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
)

type PayoutRequest struct {
	SourceAccountID string `json:"sourceAccountId"`
	Destination     struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"destination"`
	Currency          string     `json:"currency"`
	Amount            *big.Float `json:"amount"`
	Reference         string     `json:"reference"`
	ExternalReference string     `json:"externalReference"`
}

type PayoutResponse struct {
	ID                string `json:"id"`
	Status            string `json:"status"`
	CreatedDate       string `json:"createdDate"`
	ExternalReference string `json:"externalReference"`
	ApprovalStatus    string `json:"approvalStatus"`
	Message           string `json:"message"`
}

func (c *Client) InitiatePayout(payoutRequest *PayoutRequest) (*PayoutResponse, error) {
	body, err := json.Marshal(payoutRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(c.buildEndpoint("payments"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res PayoutResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetPayout(payoutID string) (*PayoutResponse, error) {
	resp, err := c.httpClient.Get(c.buildEndpoint("payments?id=%s", payoutID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res PayoutResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

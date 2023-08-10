package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Beneficiary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created string `json:"created"`
}

func (m *Client) GetBeneficiaries() ([]*Beneficiary, error) {
	resp, err := m.httpClient.Get(m.buildEndpoint("beneficiaries"))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res responseWrapper[[]*Beneficiary]
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Content, nil
}

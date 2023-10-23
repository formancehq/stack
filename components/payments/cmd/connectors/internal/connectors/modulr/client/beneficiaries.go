package client

import (
	"encoding/json"
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
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res responseWrapper[[]*Beneficiary]
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Content, nil
}

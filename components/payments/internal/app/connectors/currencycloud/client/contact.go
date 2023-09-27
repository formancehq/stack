package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Contact struct {
	ID string `json:"id"`
}

func (c *Client) GetContactID(ctx context.Context, accountID string) (*Contact, error) {
	form := url.Values{}
	form.Set("account_id", accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.buildEndpoint("/v2/contacts/find"), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	type Contacts struct {
		Contacts []*Contact `json:"contacts"`
	}

	var res Contacts
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if len(res.Contacts) == 0 {
		return nil, fmt.Errorf("no contact found for account %s", accountID)
	}

	return res.Contacts[0], nil
}

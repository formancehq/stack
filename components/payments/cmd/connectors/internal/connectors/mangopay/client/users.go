package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type user struct {
	ID string `json:"Id"`
}

func (c *Client) GetAllUsers(ctx context.Context) ([]*user, error) {
	var users []*user

	for page := 1; ; page++ {
		pagedUsers, err := c.getUsers(ctx, page)
		if err != nil {
			return nil, err
		}

		if len(pagedUsers) == 0 {
			break
		}

		users = append(users, pagedUsers...)
	}

	return users, nil
}

func (c *Client) getUsers(ctx context.Context, page int) ([]*user, error) {
	endpoint := fmt.Sprintf("%s/v2.01/%s/users", c.endpoint, c.clientID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	q := req.URL.Query()
	q.Add("per_page", "100")
	q.Add("page", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var users []*user
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users response body: %w", err)
	}

	return users, nil
}

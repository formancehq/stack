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

type user struct {
	ID           string `json:"Id"`
	CreationDate int64  `json:"CreationDate"`
}

func (c *Client) GetAllUsers(ctx context.Context, lastPage int, pageSize int) ([]*user, int, error) {
	var users []*user
	var currentPage int

	for currentPage = lastPage; ; currentPage++ {
		pagedUsers, err := c.getUsers(ctx, currentPage, pageSize)
		if err != nil {
			return nil, lastPage, err
		}

		if len(pagedUsers) == 0 {
			break
		}

		users = append(users, pagedUsers...)

		if len(pagedUsers) < pageSize {
			break
		}
	}

	return users, currentPage, nil
}

func (c *Client) getUsers(ctx context.Context, page int, pageSize int) ([]*user, error) {
	f := connectors.ClientMetrics(ctx, "mangopay", "list_users")
	now := time.Now()
	defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/users", c.endpoint, c.clientID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	q := req.URL.Query()
	q.Add("per_page", strconv.Itoa(pageSize))
	q.Add("page", fmt.Sprint(page))
	q.Add("Sort", "CreationDate:ASC")
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

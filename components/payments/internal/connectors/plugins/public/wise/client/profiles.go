package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

type Profile struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

func (c *Client) GetProfiles(ctx context.Context) ([]Profile, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "wise", "list_profiles")
	// now := time.Now()
	// defer f(ctx, now)

	var profiles []Profile
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint("v2/profiles"), http.NoBody)
	if err != nil {
		return profiles, err
	}

	var errRes wiseErrors
	statusCode, err := c.httpClient.Do(req, &profiles, &errRes)
	switch err {
	case nil:
		return profiles, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return profiles, errRes.Error(statusCode).Error()
	}
	return profiles, fmt.Errorf("failed to get profiles: %w", err)
}

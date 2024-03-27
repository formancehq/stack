package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Profile struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

func (w *Client) GetProfiles(ctx context.Context) ([]Profile, error) {
	f := connectors.ClientMetrics(ctx, "wise", "list_profiles")
	now := time.Now()
	defer f(ctx, now)

	var profiles []Profile

	res, err := w.httpClient.Get(w.endpoint("v2/profiles"))
	if err != nil {
		return profiles, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, unmarshalError(res.StatusCode, res.Body).Error()
	}

	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profiles: %w", err)
	}

	return profiles, nil
}

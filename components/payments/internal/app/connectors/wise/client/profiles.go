package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Profile struct {
	ID   uint64 `json:"id"`
	Type string `json:"type"`
}

func (w *Client) GetProfiles() ([]Profile, error) {
	var profiles []Profile

	res, err := w.httpClient.Get(w.endpoint("v1/profiles"))
	if err != nil {
		return profiles, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	err = json.Unmarshal(body, &profiles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profiles: %w", err)
	}

	return profiles, nil
}

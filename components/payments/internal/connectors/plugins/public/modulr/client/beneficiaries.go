package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

type Beneficiary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created string `json:"created"`
}

func (c *Client) GetBeneficiaries(ctx context.Context, page, pageSize int, modifiedSince time.Time) (*responseWrapper[[]Beneficiary], error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "modulr", "list_beneficiaries")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildEndpoint("beneficiaries"), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(pageSize))
	if !modifiedSince.IsZero() {
		q.Add("modifiedSince", modifiedSince.Format("2006-01-02T15:04:05-0700"))
	}
	req.URL.RawQuery = q.Encode()

	var res responseWrapper[[]Beneficiary]
	var errRes modulrError
	_, err = c.httpClient.Do(req, &res, &errRes)
	switch err {
	case nil:
		return &res, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to get beneficiaries %w", err)
}

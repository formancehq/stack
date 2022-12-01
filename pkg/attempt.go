package webhooks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/formancehq/webhooks/pkg/security"
	"github.com/pkg/errors"
)

const (
	StatusAttemptSuccess = "success"
	StatusAttemptToRetry = "to retry"
	StatusAttemptFailed  = "failed"

	KeyWebhookID      = "webhookID"
	KeyStatus         = "status"
	KeyNextRetryAfter = "nextRetryAfter"
)

type Attempt struct {
	WebhookID      string    `json:"webhookID" bson:"webhookID"`
	Date           time.Time `json:"date" bson:"date"`
	Config         Config    `json:"config" bson:"config"`
	Payload        string    `json:"payload" bson:"payload"`
	StatusCode     int       `json:"statusCode" bson:"statusCode"`
	RetryAttempt   int       `json:"retryAttempt" bson:"retryAttempt"`
	Status         string    `json:"status" bson:"status"`
	NextRetryAfter time.Time `json:"nextRetryAfter,omitempty" bson:"nextRetryAfter,omitempty"`
}

func MakeAttempt(ctx context.Context, httpClient *http.Client, schedule []time.Duration, id string, attemptNb int, cfg Config, payload []byte, isTest bool) (Attempt, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return Attempt{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	ts := time.Now().UTC()
	timestamp := ts.Unix()
	signature, err := security.Sign(id, timestamp, cfg.Secret, payload)
	if err != nil {
		return Attempt{}, errors.Wrap(err, "security.Sign")
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", "formance-webhooks/v0")
	req.Header.Set("formance-webhook-id", id)
	req.Header.Set("formance-webhook-timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("formance-webhook-signature", signature)
	req.Header.Set("formance-webhook-test", fmt.Sprintf("%v", isTest))

	resp, err := httpClient.Do(req)
	if err != nil {
		return Attempt{}, errors.Wrap(err, "http.Client.Do")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			sharedlogging.GetLogger(ctx).Error(
				errors.Wrap(err, "http.Response.Body.Close"))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Attempt{}, errors.Wrap(err, "io.ReadAll")
	}
	sharedlogging.GetLogger(ctx).Debugf("webhooks.MakeAttempt: server response body: %s", string(body))

	attempt := Attempt{
		WebhookID:    id,
		Date:         ts,
		Config:       cfg,
		Payload:      string(payload),
		StatusCode:   resp.StatusCode,
		RetryAttempt: attemptNb,
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		attempt.Status = StatusAttemptSuccess
		return attempt, nil
	}

	if attemptNb == len(schedule) {
		attempt.Status = StatusAttemptFailed
		return attempt, nil
	}

	attempt.Status = StatusAttemptToRetry
	attempt.NextRetryAfter = ts.Add(schedule[attemptNb])
	return attempt, nil
}

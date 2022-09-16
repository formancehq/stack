package webhooks

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/formancehq/webhooks/pkg/security"
	"github.com/numary/go-libs/sharedlogging"
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

func MakeAttempt(ctx context.Context, httpClient *http.Client, schedule []time.Duration, webhookID string, attemptNb int, cfg Config, data []byte) (Attempt, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.Endpoint, bytes.NewBuffer(data))
	if err != nil {
		return Attempt{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	date := time.Now().UTC()
	signature, err := security.Sign(webhookID, date, cfg.Secret, data)
	if err != nil {
		return Attempt{}, errors.Wrap(err, "security.Sign")
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", "formance-webhooks/v0")
	req.Header.Set("formance-webhook-id", webhookID)
	req.Header.Set("formance-webhook-timestamp", fmt.Sprintf("%d", date.Unix()))
	req.Header.Set("formance-webhook-signature", signature)

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

	body, _ := io.ReadAll(resp.Body)
	sharedlogging.GetLogger(ctx).Debugf("webhooks.MakeAttempt: server response body: %s\n", body)

	attempt := Attempt{
		WebhookID:    webhookID,
		Date:         date,
		Config:       cfg,
		Payload:      string(data),
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
	attempt.NextRetryAfter = date.Add(schedule[attemptNb])
	return attempt, nil
}

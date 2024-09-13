package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type webhookSubscription struct {
	Name      string `json:"name"`
	TriggerOn string `json:"trigger_on"`
	Delivery  struct {
		URL string `json:"url"`
	} `json:"delivery"`
}

func (w *Client) CreateWebhook(ctx context.Context, profileID uint64, name, triggerOn, url string) error {
	req, err := json.Marshal(webhookSubscription{
		Name:      name,
		TriggerOn: triggerOn,
		Delivery: struct {
			URL string "json:\"url\""
		}{
			URL: url,
		},
	})
	if err != nil {
		return err
	}

	res, err := w.httpClient.Post(
		w.endpoint(fmt.Sprintf("/v3/profiles/%d/subscriptions", profileID)),
		"application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

//	{
//		"data": {
//		  "resource": {
//			"type": "transfer",
//			"id": 111,
//			"profile_id": 222,
//			"account_id": 333
//		  },
//		  "current_state": "processing",
//		  "previous_state": "incoming_payment_waiting",
//		  "occurred_at": "2020-01-01T12:34:56Z"
//		},
//		"subscription_id": "01234567-89ab-cdef-0123-456789abcdef",
//		"event_type": "transfers#state-change",
//		"schema_version": "2.0.0",
//		"sent_at": "2020-01-01T12:34:56Z"
//	  }
type transferStateChangedWebhookPayload struct {
	Data struct {
		Resource struct {
			Type      string `json:"type"`
			ID        uint64 `json:"id"`
			ProfileID uint64 `json:"profile_id"`
			AccountID uint64 `json:"account_id"`
		} `json:"resource"`
		CurrentState  string `json:"current_state"`
		PreviousState string `json:"previous_state"`
		OccurredAt    string `json:"occurred_at"`
	} `json:"data"`
	SubscriptionID string `json:"subscription_id"`
	EventType      string `json:"event_type"`
	SchemaVersion  string `json:"schema_version"`
	SentAt         string `json:"sent_at"`
}

func (w *Client) TranslateTransferStateChangedWebhook(ctx context.Context, payload []byte) (Transfer, error) {
	var transferStatedChangedEvent transferStateChangedWebhookPayload
	err := json.Unmarshal(payload, &transferStatedChangedEvent)
	if err != nil {
		return Transfer{}, err
	}

	transfer, err := w.GetTransfer(ctx, fmt.Sprint(transferStatedChangedEvent.Data.Resource.ID))
	if err != nil {
		return Transfer{}, err
	}

	transfer.Created = transferStatedChangedEvent.Data.OccurredAt
	transfer.CreatedAt, err = time.Parse("2006-01-02 15:04:05", transfer.Created)
	if err != nil {
		return Transfer{}, fmt.Errorf("failed to parse created time: %w", err)
	}

	return *transfer, nil
}

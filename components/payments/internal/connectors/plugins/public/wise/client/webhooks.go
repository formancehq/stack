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
		Version string `json:"version"`
		URL     string `json:"url"`
	} `json:"delivery"`
}

type webhookSubscriptionResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Delivery struct {
		Version string `json:"version"`
		URL     string `json:"url"`
	} `json:"delivery"`
	TriggerOn string `json:"trigger_on"`
	Scope     struct {
		Domain string `json:"domain"`
	} `json:"scope"`
	CreatedBy struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"created_by"`
	CreatedAt string `json:"created_at"`
}

func (w *Client) CreateWebhook(ctx context.Context, profileID uint64, name, triggerOn, url, version string) (*webhookSubscriptionResponse, error) {
	req, err := json.Marshal(webhookSubscription{
		Name:      name,
		TriggerOn: triggerOn,
		Delivery: struct {
			Version string `json:"version"`
			URL     string `json:"url"`
		}{
			Version: version,
			URL:     url,
		},
	})
	if err != nil {
		return nil, err
	}

	res, err := w.httpClient.Post(
		w.endpoint(fmt.Sprintf("/v3/profiles/%d/subscriptions", profileID)),
		"application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response webhookSubscriptionResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (w *Client) ListWebhooksSubscription(ctx context.Context, profileID uint64) ([]webhookSubscriptionResponse, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, w.endpoint(fmt.Sprintf("/v3/profiles/%d/subscriptions", profileID)), http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response []webhookSubscriptionResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

func (w *Client) DeleteWebhooks(ctx context.Context, profileID uint64, subscriptionID string) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodDelete, w.endpoint(fmt.Sprintf("/v3/profiles/%d/subscriptions/%s", profileID, subscriptionID)), http.NoBody)
	if err != nil {
		return err
	}

	res, err := w.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

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

type balanceUpdateWebhookPayload struct {
	Data struct {
		Resource struct {
			ID        uint64 `json:"id"`
			ProfileID uint64 `json:"profile_id"`
			Type      string `json:"type"`
		} `json:"resource"`
		Amount            json.Number `json:"amount"`
		BalanceID         uint64      `json:"balance_id"`
		Currency          string      `json:"currency"`
		TransactionType   string      `json:"transaction_type"`
		OccurredAt        string      `json:"occurred_at"`
		TransferReference string      `json:"transfer_reference"`
		ChannelName       string      `json:"channel_name"`
	} `json:"data"`
	SubscriptionID string `json:"subscription_id"`
	EventType      string `json:"event_type"`
	SchemaVersion  string `json:"schema_version"`
	SentAt         string `json:"sent_at"`
}

func (w *Client) TranslateBalanceUpdateWebhook(ctx context.Context, payload []byte) (balanceUpdateWebhookPayload, error) {
	var balanceUpdateEvent balanceUpdateWebhookPayload
	err := json.Unmarshal(payload, &balanceUpdateEvent)
	if err != nil {
		return balanceUpdateWebhookPayload{}, err
	}

	return balanceUpdateEvent, nil
}

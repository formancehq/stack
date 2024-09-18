package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/uptrace/bun"
)

type webhook struct {
	bun.BaseModel `bun:"table:webhooks"`

	// Mandatory fields
	ID          string             `bun:"id,pk,type:uuid,notnull"`
	ConnectorID models.ConnectorID `bun:"connector_id,type:character varying,notnull"`

	// Optional fields
	Headers     map[string][]string `bun:"headers,type:json"`
	QueryValues map[string][]string `bun:"query_values,type:json"`
	Body        []byte              `bun:"body,type:bytea,nullzero"`
}

func (s *store) WebhooksInsert(ctx context.Context, webhook models.Webhook) error {
	toInsert := fromWebhookModels(webhook)

	_, err := s.db.NewInsert().
		Model(&toInsert).
		On("CONFLICT (id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("insert webhook", err)
	}

	return nil
}

func (s *store) WebhooksDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*webhook)(nil)).
		Where("connector_id = ?", connectorID).
		Exec(ctx)
	if err != nil {
		return e("delete webhook", err)
	}

	return nil
}

func fromWebhookModels(from models.Webhook) webhook {
	return webhook{
		ID:          from.ID,
		ConnectorID: from.ConnectorID,
		Headers:     from.Headers,
		QueryValues: from.QueryValues,
		Body:        from.Body,
	}
}

package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *Storage) CreateWebhook(ctx context.Context, webhook *models.Webhook) error {
	_, err := s.db.NewInsert().Model(webhook).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateWebhookRequestBody(ctx context.Context, webhookID uuid.UUID, requestBody []byte) error {
	if len(requestBody) == 0 {
		return errors.New("requestBody cannot be empty")
	}

	_, err := s.db.NewUpdate().Model((*models.Webhook)(nil)).Set("request_body = ?", requestBody).Where("id = ?", webhookID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetWebhook(ctx context.Context, id uuid.UUID) (*models.Webhook, error) {
	webhook := &models.Webhook{}
	err := s.db.NewSelect().Model(webhook).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return webhook, nil
}

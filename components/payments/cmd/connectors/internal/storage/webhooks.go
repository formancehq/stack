package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Storage) CreateWebhook(ctx context.Context, webhook *models.Webhook) error {
	_, err := s.db.NewInsert().Model(webhook).Exec(ctx)
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

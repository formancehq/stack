package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) ConnectorsHandleWebhooks(
	ctx context.Context,
	urlPath string,
	webhook models.Webhook,
) error {
	return s.engine.HandleWebhook(ctx, urlPath, webhook)
}

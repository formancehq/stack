package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) PaymentsUpdateMetadata(ctx context.Context, id models.PaymentID, metadata map[string]string) error {
	return s.storage.PaymentsUpdateMetadata(ctx, id, metadata)
}

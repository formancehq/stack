package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
)

func (s *Service) PaymentsGet(ctx context.Context, id models.PaymentID) (*models.Payment, error) {
	return s.storage.PaymentsGet(ctx, id)
}

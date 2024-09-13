package services

import (
	"context"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
)

func (s *Service) PaymentsList(ctx context.Context, query storage.ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error) {
	return s.storage.PaymentsList(ctx, query)
}

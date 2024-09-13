package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) PaymentsList(ctx context.Context, query storage.ListPaymentsQuery) (*bunpaginate.Cursor[models.Payment], error) {
	return s.storage.PaymentsList(ctx, query)
}

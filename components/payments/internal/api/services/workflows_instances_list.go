package services

import (
	"context"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
)

func (s *Service) WorkflowsInstancesList(ctx context.Context, query storage.ListInstancesQuery) (*bunpaginate.Cursor[models.Instance], error) {
	cursor, err := s.storage.InstancesList(ctx, query)
	return cursor, newStorageError(err, "failed to list instances")
}

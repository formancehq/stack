package services

import (
	"context"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
)

func (s *Service) ConnectorsList(ctx context.Context, query storage.ListConnectorsQuery) (*bunpaginate.Cursor[models.Connector], error) {
	cursor, err := s.storage.ConnectorsList(ctx, query)
	return cursor, newStorageError(err, "failed to list connectors")
}

package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
)

func (s *Service) ConnectorsList(ctx context.Context, query storage.ListConnectorsQuery) (*bunpaginate.Cursor[models.Connector], error) {
	cursor, err := s.storage.ConnectorsList(ctx, query)
	return cursor, newStorageError(err, "failed to list connectors")
}

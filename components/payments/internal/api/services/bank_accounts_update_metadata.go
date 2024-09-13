package services

import (
	"context"

	"github.com/google/uuid"
)

func (s *Service) BankAccountsUpdateMetadata(ctx context.Context, id uuid.UUID, metadata map[string]string) error {
	return s.storage.BankAccountsUpdateMetadata(ctx, id, metadata)
}

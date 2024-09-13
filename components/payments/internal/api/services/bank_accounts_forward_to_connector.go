package services

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
)

func (s *Service) BankAccountsForwardToConnector(ctx context.Context, bankAccountID uuid.UUID, connectorID models.ConnectorID) (*models.BankAccount, error) {
	return s.engine.ForwardBankAccount(ctx, bankAccountID, connectorID)
}

package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type ReverseTransferInitiationRequest struct {
	Reference   string            `json:"reference"`
	Description string            `json:"description"`
	Amount      *big.Int          `json:"amount"`
	Asset       string            `json:"asset"`
	Metadata    map[string]string `json:"metadata"`
}

func (r *ReverseTransferInitiationRequest) Validate() error {
	if r.Reference == "" {
		return errors.New("reference is required")
	}

	if r.Amount == nil {
		return errors.New("amount is required")
	}

	if r.Asset == "" {
		return errors.New("asset is required")
	}

	return nil
}

func checkIfReversalIsValid(transfer *models.TransferInitiation, req *ReverseTransferInitiationRequest) error {
	finalAmount := new(big.Int)
	finalAmount.Sub(transfer.Amount, req.Amount)
	switch finalAmount.Cmp(big.NewInt(0)) {
	case 0, 1:
		// Nothing to do, requested reversed amount if less than or equal to the transfer amount
	case -1:
		return errors.New("reversed amount is greater than the transfer amount")
	}

	if transfer.Type == models.TransferInitiationTypePayout {
		return errors.New("payouts cannot be reversed")
	}

	foundProcessed := false
	for _, adjustment := range transfer.RelatedAdjustments {
		if adjustment.Status == models.TransferInitiationStatusProcessed {
			foundProcessed = true
			break
		}
	}

	if !foundProcessed {
		// transfer was never processed, so we can't reverse it
		return errors.New("transfer was never processed")
	}

	return nil
}

func (s *Service) ReverseTransferInitiation(ctx context.Context, transferID string, req *ReverseTransferInitiationRequest) (*models.TransferReversal, error) {
	transferInitiationID, err := models.TransferInitiationIDFromString(transferID)
	if err != nil {
		return nil, ErrInvalidID
	}

	transfer, err := s.store.ReadTransferInitiation(ctx, transferInitiationID)
	if err != nil {
		return nil, newStorageError(err, "fetching transfer initiation")
	}

	if err := checkIfReversalIsValid(transfer, req); err != nil {
		return nil, errors.Wrap(ErrValidation, err.Error())
	}

	now := time.Now().UTC()
	reversal := &models.TransferReversal{
		ID: models.TransferReversalID{
			Reference:   req.Reference,
			ConnectorID: transfer.ConnectorID,
		},
		TransferInitiationID: transferInitiationID,
		CreatedAt:            now,
		UpdatedAt:            now,
		Description:          req.Description,
		ConnectorID:          transfer.ConnectorID,
		Amount:               req.Amount,
		Asset:                models.Asset(req.Asset),
		Status:               models.TransferReversalStatusProcessing,
		Error:                "",
		Metadata:             req.Metadata,
	}

	if err := s.store.CreateTransferReversal(ctx, reversal); err != nil {
		return nil, newStorageError(err, "creating transfer reversal")
	}

	handlers, ok := s.connectorHandlers[transfer.Provider]
	if !ok {
		return nil, errors.Wrap(ErrValidation, fmt.Sprintf("no reverse payment handler for provider %v", transfer.Provider))
	}

	if err := handlers.ReversePaymentHandler(ctx, reversal); err != nil {
		switch {
		case errors.Is(err, manager.ErrValidation):
			return nil, errors.Wrap(ErrValidation, err.Error())
		case errors.Is(err, manager.ErrConnectorNotFound):
			return nil, errors.Wrap(ErrValidation, err.Error())
		default:
			return nil, err
		}
	}

	return nil, nil
}

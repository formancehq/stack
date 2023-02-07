package stripe

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/app/models"

	"github.com/stripe/stripe-go/v72/transfer"

	"github.com/formancehq/go-libs/logging"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"

	"github.com/formancehq/payments/internal/app/ingestion"
)

func TransferTask(config Config, transferID uuid.UUID) func(ctx context.Context, ingester ingestion.Ingester, logger logging.Logger) error {
	return func(ctx context.Context, ingester ingestion.Ingester, logger logging.Logger) error {
		transferToExecute, err := ingester.GetTransfer(ctx, transferID)
		if err != nil {
			return fmt.Errorf("failed to get transfer: %w", err)
		}

		stripe.Key = config.APIKey

		params := &stripe.TransferParams{
			Params: stripe.Params{
				Context: ctx,
			},
			Amount:      stripe.Int64(transferToExecute.Amount),
			Currency:    stripe.String(transferToExecute.Currency),
			Destination: stripe.String(transferToExecute.Destination),
		}

		var (
			transferError  string
			transferStatus = models.TransferStatusSucceeded
		)

		transferResponse, err := transfer.New(params)
		if err != nil {
			transferError = err.Error()
			logger.Errorf("failed to create transfer (%s): %v", transferID, err)
		}

		if transferError != "" {
			transferStatus = models.TransferStatusFailed
		}

		err = ingester.UpdateTransferStatus(ctx, transferID, transferStatus, transferResponse.BalanceTransaction.ID, transferError)
		if err != nil {
			return fmt.Errorf("failed to update transfer status: %w", err)
		}

		return nil
	}
}

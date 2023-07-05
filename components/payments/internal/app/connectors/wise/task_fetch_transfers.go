package wise

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/wise/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func taskFetchTransfers(logger logging.Logger, c *client.Client, profileID uint64) task.Task {
	return func(
		ctx context.Context,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
	) error {
		transfers, err := c.GetTransfers(ctx, &client.Profile{
			ID: profileID,
		})
		if err != nil {
			return err
		}

		if len(transfers) == 0 {
			logger.Info("No transfers found")

			return nil
		}

		var (
			accountBatch ingestion.AccountBatch
			paymentBatch ingestion.PaymentBatch
		)

		for _, transfer := range transfers {
			logger.Info(transfer)

			var rawData json.RawMessage

			rawData, err = json.Marshal(transfer)
			if err != nil {
				return fmt.Errorf("failed to marshal transfer: %w", err)
			}

			batchElement := ingestion.PaymentBatchElement{
				Payment: &models.Payment{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: fmt.Sprintf("%d", transfer.ID),
							Type:      models.PaymentTypeTransfer,
						},
						Provider: models.ConnectorProviderWise,
					},
					CreatedAt: transfer.CreatedAt,
					Reference: fmt.Sprintf("%d", transfer.ID),
					Type:      models.PaymentTypeTransfer,
					Status:    matchTransferStatus(transfer.Status),
					Scheme:    models.PaymentSchemeOther,
					Amount:    int64(transfer.TargetValue * 100),
					Asset:     models.PaymentAsset(fmt.Sprintf("%s/2", transfer.TargetCurrency)),
					RawData:   rawData,
				},
			}

			if transfer.SourceAccount != 0 {
				ref := fmt.Sprintf("%d", transfer.SourceAccount)

				accountBatch = append(accountBatch,
					ingestion.AccountBatchElement{
						Reference: ref,
						Type:      models.AccountTypeSource,
					},
				)

				batchElement.Payment.Account = &models.Account{Reference: ref}
			}

			if transfer.TargetAccount != 0 {
				ref := fmt.Sprintf("%d", transfer.TargetAccount)

				accountBatch = append(accountBatch,
					ingestion.AccountBatchElement{
						Reference: ref,
						Provider:  models.ConnectorProviderWise.String(),
						Type:      models.AccountTypeTarget,
					},
				)

				batchElement.Payment.Account = &models.Account{Reference: ref}
			}

			paymentBatch = append(paymentBatch, batchElement)
		}

		if len(accountBatch) > 0 {
			err = ingester.IngestAccounts(ctx, accountBatch)
			if err != nil {
				return err
			}
		}

		err = ingester.IngestPayments(ctx, paymentBatch, struct{}{})
		if err != nil {
			return err
		}

		// TODO: Implement proper looper & abstract the logic

		time.Sleep(time.Minute)

		descriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{
			Name: "Fetch profiles from client",
			Key:  taskNameFetchProfiles,
		})
		if err != nil {
			return err
		}

		return scheduler.Schedule(ctx, descriptor, models.TaskSchedulerOptions{
			ScheduleOption: models.OPTIONS_RUN_NOW,
			Restart:        true,
		})
	}
}

func matchTransferStatus(status string) models.PaymentStatus {
	switch status {
	case "incoming_payment_waiting", "processing":
		return models.PaymentStatusPending
	case "funds_converted", "outgoing_payment_sent":
		return models.PaymentStatusSucceeded
	case "bounced_back", "funds_refunded":
		return models.PaymentStatusFailed
	case "cancelled":
		return models.PaymentStatusCancelled
	}

	return models.PaymentStatusOther
}

package storage

import (
	"context"
	"time"

	"github.com/formancehq/payments/internal/models"
)

func (s *Storage) UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error {
	var metadataToInsert []models.PaymentMetadata // nolint:prealloc // it's against a map

	for key, value := range metadata {
		metadataToInsert = append(metadataToInsert, models.PaymentMetadata{
			PaymentID: paymentID,
			Key:       key,
			Value:     value,
			Changelog: []models.MetadataChangelog{
				{
					CreatedAt: time.Now(),
					Value:     value,
				},
			},
		})
	}

	_, err := s.db.NewInsert().
		Model(&metadataToInsert).
		On("CONFLICT (payment_id, key) DO UPDATE").
		Set("value = EXCLUDED.value").
		Set("changelog = metadata.changelog || EXCLUDED.changelog").
		Where("metadata.value != EXCLUDED.value").
		Exec(ctx)
	if err != nil {
		return e("failed to update payment metadata", err)
	}

	return nil
}

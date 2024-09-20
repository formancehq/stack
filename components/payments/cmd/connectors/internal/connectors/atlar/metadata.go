package atlar

import (
	"fmt"
	"time"

	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/payments/internal/models"
)

const (
	atlarMetadataSpecNamespace = "com.atlar.spec/"
	valueTRUE                  = "TRUE"
	valueFALSE                 = "FALSE"
)

func ComputeAccountMetadata(key, value string) metadata.Metadata {
	namespacedKey := fmt.Sprintf("%s%s", atlarMetadataSpecNamespace, key)
	return metadata.Metadata{
		namespacedKey: value,
	}
}

func ComputeAccountMetadataBool(key string, value bool) metadata.Metadata {
	computedValue := valueFALSE
	if value {
		computedValue = valueTRUE
	}
	return ComputeAccountMetadata(key, computedValue)
}

func ComputePaymentMetadata(paymentId models.PaymentID, key, value string) *models.PaymentMetadata {
	namespacedKey := fmt.Sprintf("%s%s", atlarMetadataSpecNamespace, key)
	return &models.PaymentMetadata{
		PaymentID: paymentId,
		CreatedAt: time.Now(),
		Key:       namespacedKey,
		Value:     value,
	}
}

func ComputePaymentMetadataBool(paymentId models.PaymentID, key string, value bool) *models.PaymentMetadata {
	computedValue := valueFALSE
	if value {
		computedValue = valueTRUE
	}
	return ComputePaymentMetadata(paymentId, key, computedValue)
}

func ExtractNamespacedMetadata(metadata map[string]string, key string) (*string, error) {
	value, ok := metadata[atlarMetadataSpecNamespace+key]
	if !ok {
		return nil, fmt.Errorf("unable to find metadata with key %s%s", atlarMetadataSpecNamespace, key)
	}
	return &value, nil
}

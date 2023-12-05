package atlar

import (
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

const (
	formanceNamespace = "com.atlar.spec/"
	valueTRUE         = "TRUE"
	valueFALSE        = "FALSE"
)

func ComputeAccountMetadata(key, value string) metadata.Metadata {
	namespacedKey := fmt.Sprintf("%s%s", formanceNamespace, key)
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

func ComputePaymentMetadata(paymentId models.PaymentID, key, value string) *models.Metadata {
	namespacedKey := fmt.Sprintf("%s%s", formanceNamespace, key)
	return &models.Metadata{
		PaymentID: paymentId,
		CreatedAt: time.Now(),
		Key:       namespacedKey,
		Value:     value,
	}
}

func ComputePaymentMetadataBool(paymentId models.PaymentID, key string, value bool) *models.Metadata {
	computedValue := valueFALSE
	if value {
		computedValue = valueTRUE
	}
	return ComputePaymentMetadata(paymentId, key, computedValue)
}

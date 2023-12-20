package dummypay

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
)

const (
	taskKeyIngest = "ingest"
)

// newTaskIngest returns a new task descriptor for the ingest task.
func newTaskIngest(filePath string) TaskDescriptor {
	return TaskDescriptor{
		Name:     "Ingest accounts and payments from read files",
		Key:      taskKeyIngest,
		FileName: filePath,
	}
}

// taskIngest ingests a payment file.
func taskIngest(config Config, descriptor TaskDescriptor, fs fs) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		err := handleFile(ctx, connectorID, ingester, config, descriptor, fs)
		if err != nil {
			return fmt.Errorf("failed to handle file '%s': %w", descriptor.FileName, err)
		}

		return nil
	}
}

func handleFile(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	config Config,
	descriptor TaskDescriptor,
	fs fs,
) error {
	object, err := getObject(config, descriptor, fs)
	if err != nil {
		return err
	}

	switch object.Kind {
	case KindAccount:
		batch, err := handleAccount(connectorID, object.Account)
		if err != nil {
			return err
		}

		err = ingester.IngestAccounts(ctx, batch)
		if err != nil {
			return fmt.Errorf("failed to ingest file '%s': %w", descriptor.FileName, err)
		}
	case KindPayment:
		batch, err := handlePayment(connectorID, object.Payment)
		if err != nil {
			return err
		}

		err = ingester.IngestPayments(ctx, connectorID, batch, struct{}{})
		if err != nil {
			return fmt.Errorf("failed to ingest file '%s': %w", descriptor.FileName, err)
		}
	default:
		return fmt.Errorf("unknown object kind '%s'", object.Kind)
	}

	return nil
}

func getObject(config Config, descriptor TaskDescriptor, fs fs) (*object, error) {
	file, err := fs.Open(filepath.Join(config.Directory, descriptor.FileName))
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %w", descriptor.FileName, err)
	}
	defer file.Close()

	var objectElement object
	err = json.NewDecoder(file).Decode(&objectElement)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file '%s': %w", descriptor.FileName, err)
	}

	return &objectElement, nil
}

func handleAccount(connectorID models.ConnectorID, accountElement *account) (ingestion.AccountBatch, error) {
	accountType, err := models.AccountTypeFromString(accountElement.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account type: %w", err)
	}

	raw, err := json.Marshal(accountElement)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment: %w", err)
	}

	ingestionPayload := ingestion.AccountBatch{&models.Account{
		ID: models.AccountID{
			Reference:   accountElement.Reference,
			ConnectorID: connectorID,
		},
		ConnectorID:  connectorID,
		CreatedAt:    accountElement.CreatedAt,
		Reference:    accountElement.Reference,
		DefaultAsset: models.Asset(accountElement.DefaultAsset),
		AccountName:  accountElement.AccountName,
		Type:         accountType,
		Metadata:     map[string]string{},
		RawData:      raw,
	}}

	return ingestionPayload, nil
}

func handlePayment(connectorID models.ConnectorID, paymentElement *payment) (ingestion.PaymentBatch, error) {
	paymentType, err := models.PaymentTypeFromString(paymentElement.Type)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment type '%s': %w", paymentElement.Type, err)
	}

	paymentStatus, err := models.PaymentStatusFromString(paymentElement.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment status '%s': %w", paymentElement.Status, err)
	}

	paymentScheme, err := models.PaymentSchemeFromString(paymentElement.Scheme)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payment scheme '%s': %w", paymentElement.Scheme, err)
	}

	raw, err := json.Marshal(paymentElement)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment: %w", err)
	}

	ingestionPayload := ingestion.PaymentBatch{ingestion.PaymentBatchElement{
		Payment: &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: paymentElement.Reference,
					Type:      paymentType,
				},
				ConnectorID: connectorID,
			},
			Reference:     paymentElement.Reference,
			ConnectorID:   connectorID,
			Amount:        paymentElement.Amount,
			InitialAmount: paymentElement.Amount,
			Type:          paymentType,
			Status:        paymentStatus,
			Scheme:        paymentScheme,
			Asset:         models.Asset(paymentElement.Asset),
			RawData:       raw,
		},
		Update: true,
	}}

	return ingestionPayload, nil
}

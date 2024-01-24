package dummypay

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
)

const (
	taskKeyInitDirectory = "init-directory"
	asset                = "DUMMYCOIN"
	generatedFilePrefix  = "dummypay-generated-file"
)

// newTaskGenerateFiles returns a new task descriptor for the task that generates files.
func newTaskGenerateFiles() TaskDescriptor {
	return TaskDescriptor{
		Name: "Generate files into a directory",
		Key:  taskKeyInitDirectory,
	}
}

// taskGenerateFiles generates payment files to a given directory.
func taskGenerateFiles(config Config, fs fs) task.Task {
	return func(ctx context.Context, ingester ingestion.Ingester, connectorID models.ConnectorID) error {
		err := fs.Mkdir(config.Directory, 0o777) //nolint:gomnd
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf(
				"failed to create dummypay config directory '%s': %w", config.Directory, err)
		}

		var accountIDs []*models.AccountID
		for i := 0; i < config.NumberOfAccountsPreGenerated; i++ {
			accountID, err := generateAccountsFile(ctx, connectorID, ingester, config, fs)
			if err != nil {
				return fmt.Errorf("failed to generate accounts file: %w", err)
			}

			accountIDs = append(accountIDs, accountID)
		}

		for i := 0; i < config.NumberOfPaymentsPreGenerated; i++ {
			if err := generatePaymentsFile(ctx, generatedFilePrefix, connectorID, ingester, accountIDs, config, fs); err != nil {
				return fmt.Errorf("failed to generate payments files: %w", err)
			}
		}

		return nil
	}
}

func generateAccountsFile(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	config Config,
	fs fs,
) (*models.AccountID, error) {
	name := fmt.Sprintf("account-%d", time.Now().UnixNano())
	key := fmt.Sprintf("%s-%s", generatedFilePrefix, name)
	fileKey := fmt.Sprintf("%s/%s.json", config.Directory, key)

	generatedAccount := account{
		Reference:    uuid.New().String(),
		CreatedAt:    time.Now(),
		DefaultAsset: asset,
		AccountName:  name,
		Type:         generateRandomAccountType().String(),
		Metadata:     map[string]string{},
	}

	file, err := fs.Create(fileKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode the payment object as JSON to a new file.
	err = json.NewEncoder(file).Encode(&object{
		Kind:    KindAccount,
		Account: &generatedAccount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode json into file: %w", err)
	}

	raw, err := json.Marshal(generatedAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal account: %w", err)
	}

	accountID := models.AccountID{
		Reference:   generatedAccount.Reference,
		ConnectorID: connectorID,
	}
	if err := ingester.IngestAccounts(ctx, ingestion.AccountBatch{
		{
			ID:           accountID,
			ConnectorID:  connectorID,
			CreatedAt:    generatedAccount.CreatedAt,
			Reference:    generatedAccount.Reference,
			DefaultAsset: asset,
			AccountName:  name,
			Type:         models.AccountType(generatedAccount.Type),
			Metadata:     map[string]string{},
			RawData:      raw,
		},
	}); err != nil {
		return nil, fmt.Errorf("failed to ingest accounts: %w", err)
	}

	return &accountID, nil
}

func generatePaymentsFile(
	ctx context.Context,
	prefix string,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	accountIDs []*models.AccountID,
	config Config,
	fs fs,
) error {
	name := fmt.Sprintf("payment-%d", time.Now().UnixNano())
	key := name
	if prefix != "" {
		key = fmt.Sprintf("%s-%s", generatedFilePrefix, name)
	}
	fileKey := fmt.Sprintf("%s/%s.json", config.Directory, key)

	generatedPayment := payment{
		Reference: uuid.New().String(),
		CreatedAt: time.Now(),
		Amount:    big.NewInt(int64(rand.Intn(10000))),
		Asset:     asset,
		Type:      generateRandomPaymentType().String(),
		Status:    generateRandomStatus().String(),
		Scheme:    generateRandomScheme().String(),
		Metadata:  map[string]string{},
	}

	var sourceAccountID, destinationAccountID *models.AccountID
	if len(accountIDs) != 0 {
		if generateRandomNumber() > nMax/2 {
			sourceAccountID = accountIDs[generateRandomNumber()%len(accountIDs)]
			generatedPayment.SourceAccountID = sourceAccountID.String()
		} else {
			destinationAccountID = accountIDs[generateRandomNumber()%len(accountIDs)]
			generatedPayment.DestinationAccountID = destinationAccountID.String()
		}
	}

	file, err := fs.Create(fileKey)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode the payment object as JSON to a new file.
	err = json.NewEncoder(file).Encode(&object{
		Kind:    KindPayment,
		Payment: &generatedPayment,
	})
	if err != nil {
		return fmt.Errorf("failed to encode json into file: %w", err)
	}

	raw, err := json.Marshal(generatedPayment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}
	if err := ingester.IngestPayments(ctx, ingestion.PaymentBatch{
		{
			Payment: &models.Payment{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: generatedPayment.Reference,
						Type:      models.PaymentType(generatedPayment.Type),
					},
					ConnectorID: connectorID,
				},
				ConnectorID:          connectorID,
				CreatedAt:            generatedPayment.CreatedAt,
				Reference:            generatedPayment.Reference,
				Amount:               generatedPayment.Amount,
				InitialAmount:        generatedPayment.Amount,
				Type:                 models.PaymentType(generatedPayment.Type),
				Status:               models.PaymentStatus(generatedPayment.Status),
				Scheme:               models.PaymentScheme(generatedPayment.Scheme),
				Asset:                asset,
				RawData:              raw,
				SourceAccountID:      sourceAccountID,
				DestinationAccountID: destinationAccountID,
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to ingest payments: %w", err)
	}

	return nil
}

// nMax is the maximum number that can be generated
// with the minimum being 0.
const nMax = 10000

func generateRandomAccountType() models.AccountType {
	// 50% chance.
	accountType := models.AccountTypeInternal

	// 50% chance.
	if generateRandomNumber() > nMax/2 {
		accountType = models.AccountTypeExternal
	}

	return accountType
}

// generateRandomNumber generates a random number between 0 and nMax.
func generateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())

	//nolint:gosec // allow weak random number generator as it is not used for security
	value := rand.Intn(nMax)

	return value
}

// generateRandomType generates a random payment type.
func generateRandomPaymentType() models.PaymentType {
	paymentType := models.PaymentTypePayIn

	num := generateRandomNumber()
	switch {
	case num < nMax/4: // 25% chance
		paymentType = models.PaymentTypePayOut
	case num < nMax/3: // ~9% chance
		paymentType = models.PaymentTypeTransfer
	}

	return paymentType
}

// generateRandomStatus generates a random payment status.
func generateRandomStatus() models.PaymentStatus {
	// ~50% chance.
	paymentStatus := models.PaymentStatusSucceeded

	num := generateRandomNumber()

	switch {
	case num < nMax/4: // 25% chance
		paymentStatus = models.PaymentStatusPending
	case num < nMax/3: // ~9% chance
		paymentStatus = models.PaymentStatusFailed
	case num < nMax/2: // ~16% chance
		paymentStatus = models.PaymentStatusCancelled
	}

	return paymentStatus
}

// generateRandomScheme generates a random payment scheme.
func generateRandomScheme() models.PaymentScheme {
	num := generateRandomNumber() / 1000 //nolint:gomnd // allow for random number

	paymentScheme := models.PaymentSchemeCardMasterCard

	availableSchemes := []models.PaymentScheme{
		models.PaymentSchemeCardMasterCard,
		models.PaymentSchemeCardVisa,
		models.PaymentSchemeCardDiscover,
		models.PaymentSchemeCardJCB,
		models.PaymentSchemeCardUnionPay,
		models.PaymentSchemeCardAmex,
		models.PaymentSchemeCardDiners,
		models.PaymentSchemeSepaDebit,
		models.PaymentSchemeSepaCredit,
		models.PaymentSchemeApplePay,
		models.PaymentSchemeGooglePay,
		models.PaymentSchemeA2A,
		models.PaymentSchemeACHDebit,
		models.PaymentSchemeACH,
		models.PaymentSchemeRTP,
	}

	if num < len(availableSchemes) {
		paymentScheme = availableSchemes[num]
	}

	return paymentScheme
}

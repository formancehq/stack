package dummypay

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/models"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/stretchr/testify/assert"
)

// Create a minimal mock for connector installation.
type (
	mockConnectorContext struct {
		ctx context.Context
	}
	mockScheduler struct{}
)

func (mcc *mockConnectorContext) Context() context.Context {
	return mcc.ctx
}

func (mcc mockScheduler) Schedule(ctx context.Context, p models.TaskDescriptor, options models.TaskSchedulerOptions) error {
	return nil
}

func (mcc *mockConnectorContext) Scheduler() task.Scheduler {
	return mockScheduler{}
}

func TestConnector(t *testing.T) {
	t.Parallel()

	config := Config{}
	logger := logging.FromContext(context.TODO())

	fileSystem := newTestFS()

	connector := newConnector(logger, config, fileSystem)

	err := connector.Install(new(mockConnectorContext))
	assert.NoErrorf(t, err, "Install() failed")

	testCases := []struct {
		key  taskKey
		task task.Task
	}{
		{taskKeyReadFiles, taskReadFiles(config, fileSystem)},
		{taskKeyInitDirectory, taskGenerateFiles(config, fileSystem)},
		{taskKeyIngest, taskIngest(config, TaskDescriptor{}, fileSystem)},
	}

	for _, testCase := range testCases {
		var taskDescriptor models.TaskDescriptor

		taskDescriptor, err = models.EncodeTaskDescriptor(TaskDescriptor{Key: testCase.key})
		assert.NoErrorf(t, err, "EncodeTaskDescriptor() failed")

		assert.EqualValues(t,
			reflect.ValueOf(testCase.task).String(),
			reflect.ValueOf(connector.Resolve(taskDescriptor)).String(),
		)
	}

	taskDescriptor, err := models.EncodeTaskDescriptor(TaskDescriptor{Key: "test"})
	assert.NoErrorf(t, err, "EncodeTaskDescriptor() failed")

	assert.EqualValues(t,
		reflect.ValueOf(func() error { return nil }).String(),
		reflect.ValueOf(connector.Resolve(taskDescriptor)).String(),
	)

	assert.NoError(t, connector.Uninstall(context.TODO()))
}

type MockIngester struct{}

func (m *MockIngester) IngestAccounts(ctx context.Context, batch ingestion.AccountBatch) error {
	return nil
}

func (m *MockIngester) IngestPayments(ctx context.Context, batch ingestion.PaymentBatch) error {
	return nil
}

func (m *MockIngester) IngestBalances(ctx context.Context, batch ingestion.BalanceBatch, checkIfAccountExists bool) error {
	return nil
}

func (m *MockIngester) UpdateTaskState(ctx context.Context, state any) error {
	return nil
}

func (m *MockIngester) UpdateTransferInitiationPayment(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, updatedAt time.Time) error {
	return nil
}

func (m *MockIngester) UpdateTransferInitiationPaymentsStatus(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, updatedAt time.Time) error {
	return nil
}

func (m *MockIngester) AddTransferInitiationPaymentID(ctx context.Context, tf *models.TransferInitiation, paymentID *models.PaymentID, updatedAt time.Time) error {
	return nil
}

func (m *MockIngester) UpdateTransferReversalStatus(ctx context.Context, transfer *models.TransferInitiation, transferReversal *models.TransferReversal) error {
	return nil
}

func (m *MockIngester) LinkBankAccountWithAccount(ctx context.Context, bankAccount *models.BankAccount, accountID *models.AccountID) error {
	return nil
}

var _ ingestion.Ingester = (*MockIngester)(nil)

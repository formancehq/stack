package suite

import (
	"math/big"
	"time"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	fakeConnectorID = "eyJQcm92aWRlciI6Ik1PRFVMUiIsIlJlZmVyZW5jZSI6IjQyNDY5MjEwLTQ5YTMtNGQ0NS04OWQ1LWVmNWI3YWI4OTUwNyJ9"
)

var _ = WithModules([]*Module{modules.Payments}, func() {
	When("trying to create accounts of an non existent connector", func() {
		BeforeEach(func() {
			_, err := Client().Payments.V1.CreateAccount(
				TestContext(),
				shared.AccountRequest{
					AccountName: ptr("test"),
					ConnectorID: fakeConnectorID, // not installed
					CreatedAt:   time.Now(),
					Reference:   "test1",
					Type:        shared.AccountTypeInternal,
				},
			)
			Expect(err).NotTo(BeNil())
			Expect(err.(*sdkerrors.PaymentsErrorResponse).ErrorCode).To(Equal(shared.PaymentsErrorsEnumValidation))
		})
		It("Should fail with a 400", func() {
		})
	})
})

var _ = WithModules([]*Module{modules.Payments}, func() {
	When("trying to create payments of an non existent connector", func() {
		BeforeEach(func() {
			_, err := Client().Payments.V1.CreatePayment(
				TestContext(),
				shared.PaymentRequest{
					Amount:      big.NewInt(100),
					Asset:       "EUR/2",
					ConnectorID: fakeConnectorID,
					CreatedAt:   time.Now(),
					Reference:   "test",
					Scheme:      shared.PaymentSchemeOther,
					Status:      shared.PaymentStatusSucceeded,
					Type:        shared.PaymentTypeTransfer,
				},
			)
			Expect(err).NotTo(BeNil())
			Expect(err.(*sdkerrors.PaymentsErrorResponse).ErrorCode).To(Equal(shared.PaymentsErrorsEnumValidation))
		})
		It("Should fail with a 400", func() {
		})
	})
})

var _ = WithModules([]*Module{modules.Payments}, func() {
	var (
		connectorID string
	)
	BeforeEach(func() {
		response, err := Client().Payments.V1.InstallConnector(
			TestContext(),
			operations.InstallConnectorRequest{
				ConnectorConfig: shared.ConnectorConfig{
					GenericConfig: &shared.GenericConfig{
						Name: "test",
					},
				},
				Connector: shared.ConnectorGeneric,
			},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(Equal(201))
		Expect(response.ConnectorResponse).ToNot(BeNil())

		connectorID = response.ConnectorResponse.Data.ConnectorID
	})
	When("creating accounts and payments", func() {
		var (
			accountIDInternal1 string
			accountIDInternal2 string
		)
		BeforeEach(func() {
			createAccountResponse, err := Client().Payments.V1.CreateAccount(
				TestContext(),
				shared.AccountRequest{
					AccountName: ptr("test 1"),
					ConnectorID: connectorID,
					CreatedAt:   time.Now(),
					Reference:   "test1",
					Type:        shared.AccountTypeInternal,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createAccountResponse.StatusCode).To(Equal(200))

			accountIDInternal1 = createAccountResponse.PaymentsAccountResponse.Data.ID

			createAccountResponse, err = Client().Payments.V1.CreateAccount(
				TestContext(),
				shared.AccountRequest{
					AccountName: ptr("test 2"),
					ConnectorID: connectorID,
					CreatedAt:   time.Now(),
					Reference:   "test2",
					Type:        shared.AccountTypeInternal,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createAccountResponse.StatusCode).To(Equal(200))

			accountIDInternal2 = createAccountResponse.PaymentsAccountResponse.Data.ID

			createPaymentResponse, err := Client().Payments.V1.CreatePayment(
				TestContext(),
				shared.PaymentRequest{
					Amount:               big.NewInt(100),
					Asset:                "EUR/2",
					ConnectorID:          connectorID,
					CreatedAt:            time.Now(),
					DestinationAccountID: &accountIDInternal2,
					Reference:            "p1",
					Scheme:               shared.PaymentSchemeOther,
					SourceAccountID:      &accountIDInternal1,
					Status:               shared.PaymentStatusSucceeded,
					Type:                 shared.PaymentTypeTransfer,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createPaymentResponse.StatusCode).To(Equal(200))

			createPaymentResponse, err = Client().Payments.V1.CreatePayment(
				TestContext(),
				shared.PaymentRequest{
					Amount:               big.NewInt(200),
					Asset:                "EUR/2",
					ConnectorID:          connectorID,
					CreatedAt:            time.Now(),
					Reference:            "p2",
					DestinationAccountID: &accountIDInternal1,
					Scheme:               shared.PaymentSchemeOther,
					Status:               shared.PaymentStatusSucceeded,
					Type:                 shared.PaymentTypePayIn,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createPaymentResponse.StatusCode).To(Equal(200))

			createPaymentResponse, err = Client().Payments.V1.CreatePayment(
				TestContext(),
				shared.PaymentRequest{
					Amount:          big.NewInt(300),
					Asset:           "EUR/2",
					ConnectorID:     connectorID,
					CreatedAt:       time.Now(),
					Reference:       "p3",
					SourceAccountID: &accountIDInternal1,
					Scheme:          shared.PaymentSchemeOther,
					Status:          shared.PaymentStatusFailed,
					Type:            shared.PaymentTypePayout,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createPaymentResponse.StatusCode).To(Equal(200))
		})
		It("should be available on api", func() {
			listAccountsResponse, err := Client().Payments.V1.PaymentslistAccounts(
				TestContext(),
				operations.PaymentslistAccountsRequest{},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(listAccountsResponse.StatusCode).To(Equal(200))
			Expect(listAccountsResponse.AccountsCursor.Cursor.Data).To(HaveLen(2))

			listPaymentsResponse, err := Client().Payments.V1.ListPayments(
				TestContext(),
				operations.ListPaymentsRequest{},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(listPaymentsResponse.StatusCode).To(Equal(200))
			Expect(listPaymentsResponse.PaymentsCursor.Cursor.Data).To(HaveLen(3))
		})
	})
})

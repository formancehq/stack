package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"os"
	"path/filepath"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	paymentEvents "github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/events"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Payments}, func() {
	When("configuring dummy pay connector", func() {
		var (
			msgs               chan *nats.Msg
			cancelSubscription func()
		)
		BeforeEach(func() {
			cancelSubscription, msgs = SubscribePayments()

			paymentsDir := filepath.Join(os.TempDir(), uuid.NewString())
			Expect(os.MkdirAll(paymentsDir, 0o777)).To(Succeed())
			response, err := Client().Payments.InstallConnector(
				TestContext(),
				operations.InstallConnectorRequest{
					ConnectorConfig: shared.ConnectorConfig{
						DummyPayConfig: &shared.DummyPayConfig{
							FilePollingPeriod:    ptr("1s"),
							Directory:            paymentsDir,
							FileGenerationPeriod: ptr("1s"),
						},
					},
					Connector: shared.ConnectorDummyPay,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should trigger some events", func() {
			msg := WaitOnChanWithTimeout(msgs, 10*time.Second)
			Expect(events.Check(msg.Data, "payments", paymentEvents.EventTypeSavedPayments)).Should(Succeed())
		})
		It("should generate some payments", func() {
			Eventually(func(g Gomega) []shared.Payment {
				response, err := Client().Payments.ListPayments(
					TestContext(),
					operations.ListPaymentsRequest{},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				return response.PaymentsCursor.Cursor.Data
			}).WithTimeout(10 * time.Second).ShouldNot(BeEmpty()) // TODO: Check other fields
		})
	})
})

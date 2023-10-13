package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"os"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	paymentEvents "github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/events"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Payments}, func() {
	When("configuring stripe connector", func() {
		var (
			msgs               chan *nats.Msg
			cancelSubscription func()
		)
		BeforeEach(func() {
			apiKey := os.Getenv("STRIPE_API_KEY")
			if apiKey == "" {
				Skip("No stripe api key provided")
			}

			cancelSubscription, msgs = SubscribePayments()

			response, err := Client().Payments.InstallConnector(
				TestContext(),
				operations.InstallConnectorRequest{
					RequestBody: shared.StripeConfig{
						APIKey: apiKey,
					},
					Connector: shared.ConnectorStripe,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should trigger some events", func() {
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
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
			}).ShouldNot(BeEmpty()) // TODO: Check other fields
		})
		It("should be ingested on search", func() {
			Eventually(func(g Gomega) bool {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("PAYMENT"),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))
				g.Expect(response.Response.Cursor.Data).NotTo(BeEmpty())

				return true
			}).Should(BeTrue())
		})
	})
})

package suite_test

import (
	"os"
	"time"

	"github.com/formancehq/formance-sdk-go"
	paymentEvents "github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/events"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
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

			_, err := Client().PaymentsApi.
				InstallConnector(TestContext(), formance.STRIPE).
				ConnectorConfig(formance.ConnectorConfig{
					StripeConfig: &formance.StripeConfig{
						ApiKey: apiKey,
					},
				}).
				Execute()
			Expect(err).To(BeNil())
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should trigger some events", func() {
			msg := WaitOnChanWithTimeout(msgs, 5*time.Second)
			Expect(events.Check(msg.Data, "payments", paymentEvents.EventTypeSavedPayments)).Should(BeNil())
		})
		It("should generate some payments", func() {
			Eventually(func(g Gomega) []formance.Payment {
				res, _, err := Client().PaymentsApi.
					ListPayments(TestContext()).
					Execute()
				g.Expect(err).To(BeNil())
				return res.Cursor.Data
			}).ShouldNot(BeEmpty()) // TODO: Check other fields
		})
		It("should be ingested on search", func() {
			Eventually(func(g Gomega) bool {
				res, _, err := Client().SearchApi.Search(TestContext()).Query(formance.Query{
					Target: formance.PtrString("PAYMENT"),
				}).Execute()
				g.Expect(err).To(BeNil())
				g.Expect(res.Cursor.Data).NotTo(BeEmpty())

				return true
			}).Should(BeTrue())
		})
	})
})

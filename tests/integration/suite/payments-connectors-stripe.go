package suite

import (
	"os"
	"time"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/events"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
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

			ev := events.Event{}
			err := proto.Unmarshal(msg.Data, &ev)
			Expect(err).To(BeNil())

			switch ev.Event.(type) {
			case *events.Event_PaymentSaved:
				// Expect(models.PaymentType(paymentSaved.PaymentSaved.Type).IsValid()).Should(BeTrue())
				// Expect(models.PaymentStatus(paymentSaved.PaymentSaved.Status).IsValid()).Should(BeTrue())
				// Expect(models.PaymentScheme(paymentSaved.PaymentSaved.Scheme).IsValid()).Should(BeTrue())
			default:
				Expect(false).Should(BeTrue(), "Unexpected event type: %T", ev.Event)
			}
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

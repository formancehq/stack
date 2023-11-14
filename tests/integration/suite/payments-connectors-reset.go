package suite

import (
	"encoding/json"
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

var _ = Given("some environment with dummy pay connector", func() {
	var (
		existingPaymentID  string
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
				RequestBody: shared.DummyPayConfig{
					FilePollingPeriod:    ptr("1s"),
					Directory:            paymentsDir,
					FileGenerationPeriod: ptr("1s"),
				},
				Connector: shared.ConnectorDummyPay,
			},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(Equal(201))

		Eventually(func(g Gomega) bool {
			response, err := Client().Search.Search(
				TestContext(),
				shared.Query{
					Target: ptr("PAYMENT"),
				},
			)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(response.StatusCode).To(Equal(200))
			if len(response.Response.Cursor.Data) == 0 {
				return false
			}
			existingPaymentID = response.Response.Cursor.Data[0]["id"].(string)
			return true
		}).Should(BeTrue())
	})
	AfterEach(func() {
		cancelSubscription()
	})
	When("resetting connector", func() {
		BeforeEach(func() {
			response, err := Client().Payments.ResetConnector(
				TestContext(),
				operations.ResetConnectorRequest{
					Connector: shared.ConnectorDummyPay,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
		})
		It("should trigger some events", func() {
			var msg *nats.Msg
			Eventually(func(g Gomega) bool {
				msg = WaitOnChanWithTimeout(msgs, 5*time.Second)
				type typedMessage struct {
					Type string `json:"type"`
				}
				tm := &typedMessage{}
				g.Expect(json.Unmarshal(msg.Data, tm)).To(Succeed())
				return tm.Type == paymentEvents.EventTypeConnectorReset
			}).Should(BeTrue())

			Expect(events.Check(msg.Data, "payments", paymentEvents.EventTypeConnectorReset)).Should(Succeed())
		})
		It("should delete payments on search service", func() {
			Eventually(func(g Gomega) []map[string]any {
				response, err := Client().Search.Search(
					TestContext(),
					shared.Query{
						Target: ptr("PAYMENT"),
						Terms:  []string{"id=" + existingPaymentID},
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))
				return response.Response.Cursor.Data
			}).Should(BeEmpty())
		})
	})
})

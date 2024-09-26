package suite

import (
	webhooks "github.com/formancehq/webhooks/pkg"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v3/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
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
		JustBeforeEach(func() {
			cancelSubscription, msgs = SubscribePayments()

			paymentsDir := filepath.Join(os.TempDir(), uuid.NewString())
			Expect(os.MkdirAll(paymentsDir, 0o777)).To(Succeed())
			response, err := Client().Payments.V1.InstallConnector(
				TestContext(),
				operations.InstallConnectorRequest{
					ConnectorConfig: shared.ConnectorConfig{
						DummyPayConfig: &shared.DummyPayConfig{
							FilePollingPeriod:            ptr("1s"),
							Directory:                    paymentsDir,
							Name:                         "test",
							NumberOfAccountsPreGenerated: ptr(int64(0)),
							NumberOfPaymentsPreGenerated: ptr(int64(1)),
						},
					},
					Connector: shared.ConnectorDummyPay,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))
			Expect(response.ConnectorResponse).ToNot(BeNil())
		})
		JustAfterEach(func() {
			cancelSubscription()
		})
		It("should trigger some events", func() {
			msg := WaitOnChanWithTimeout(msgs, 20*time.Second)
			Expect(events.Check(msg.Data, "payments", paymentEvents.EventTypeSavedPayments)).Should(Succeed())
		})
		It("should generate some payments", func() {
			Eventually(func(g Gomega) []shared.Payment {
				response, err := Client().Payments.V1.ListPayments(
					TestContext(),
					operations.ListPaymentsRequest{},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				return response.PaymentsCursor.Cursor.Data
			}).WithTimeout(10 * time.Second).ShouldNot(BeEmpty()) // TODO: Check other fields
		})
		WithModules([]*Module{modules.Webhooks}, func() {
			var (
				httpServer *httptest.Server
				called     chan []byte
				secret     = webhooks.NewSecret()
			)
			BeforeEach(func() {
				called = make(chan []byte)
				httpServer = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						defer close(called)
						data, _ := io.ReadAll(r.Body)
						called <- data
					}))
				DeferCleanup(func() {
					httpServer.Close()
				})

				response, err := Client().Webhooks.V1.InsertConfig(
					TestContext(),
					shared.ConfigUser{
						Endpoint: httpServer.URL,
						Secret:   &secret,
						EventTypes: []string{
							"payments.saved_payment",
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			})
			It("Should trigger a webhook", func() {
				Eventually(called).Should(ReceiveEvent("payments", paymentEvents.EventTypeSavedPayments))
			})
		})
	})
})

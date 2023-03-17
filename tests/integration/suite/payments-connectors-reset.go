package suite

import (
	"os"
	"path/filepath"
	"time"

	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/events"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
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
		Expect(os.MkdirAll(paymentsDir, 0777)).To(BeNil())
		_, err := Client().PaymentsApi.
			InstallConnector(TestContext(), formance.DUMMY_PAY).
			ConnectorConfig(formance.ConnectorConfig{
				DummyPayConfig: &formance.DummyPayConfig{
					FilePollingPeriod:    formance.PtrString("1s"),
					Directory:            paymentsDir,
					FileGenerationPeriod: formance.PtrString("1s"),
				},
			}).
			Execute()
		Expect(err).To(BeNil())

		Eventually(func(g Gomega) bool {
			response, _, err := Client().SearchApi.
				Search(TestContext()).
				Query(formance.Query{
					Target: formance.PtrString("PAYMENT"),
				}).
				Execute()
			g.Expect(err).To(BeNil())
			if len(response.Cursor.Data) == 0 {
				return false
			}
			existingPaymentID = response.Cursor.Data[0].(map[string]any)["id"].(string)
			return true
		}).Should(BeTrue())
	})
	AfterEach(func() {
		cancelSubscription()
	})
	When("resetting connector", func() {
		var (
			err error
		)
		BeforeEach(func() {
			_, err = Client().PaymentsApi.
				ResetConnector(TestContext(), formance.DUMMY_PAY).
				Execute()
			Expect(err).To(BeNil())
		})
		It("should trigger some events", func() {
			var msg *nats.Msg
			Eventually(func(g Gomega) bool {
				msg = WaitOnChanWithTimeout(msgs, 5*time.Second)

				ev := events.Event{}
				err := proto.Unmarshal(msg.Data, &ev)
				Expect(err).To(BeNil())

				switch resetConnector := ev.Event.(type) {
				case *events.Event_ResetConnector:
					Expect(resetConnector.ResetConnector.Provider).To(Equal(formance.DUMMY_PAY))
					return true
				default:
					return false
				}
			}).Should(BeTrue())
		})
		It("should delete payments on search service", func() {
			Eventually(func(g Gomega) []any {
				ret, _, err := Client().SearchApi.
					Search(TestContext()).
					Query(formance.Query{
						Target: formance.PtrString("PAYMENT"),
						Terms:  []string{"id=" + existingPaymentID},
					}).
					Execute()
				g.Expect(err).To(BeNil())
				return ret.Cursor.Data
			}).Should(BeEmpty())
		})
	})
})

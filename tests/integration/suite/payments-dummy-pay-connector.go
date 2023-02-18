package suite

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/numary/ledger/pkg/bus"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("configuring dummy pay connector", func() {
		var (
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
		})
		AfterEach(func() {
			cancelSubscription()
		})
		It("should trigger some events", func() {
			// Wait for created transaction event
			msg := waitOnChanWithTimeout(msgs, 5*time.Second)
			event := &bus.EventMessage{}
			Expect(json.Unmarshal(msg.Data, event)).To(BeNil())
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
	})
})

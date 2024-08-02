package suite

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/publish"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
	"time"
)

var _ = WithModules([]*Module{modules.Search}, func() {
	When("sending a v1 created transaction event in the event bus", func() {
		now := time.Now()
		BeforeEach(func() {
			CreateTopic("ledger")

			PublishLedger(publish.EventMessage{
				Date:    time.Now(),
				App:     "ledger",
				Version: "v1",
				Type:    "COMMITTED_TRANSACTIONS",
				Payload: map[string]any{
					"ledger": "foo",
					"transactions": []map[string]any{
						{
							"preCommitVolumes":  map[string]map[string]*big.Int{},
							"postCommitVolumes": map[string]map[string]*big.Int{},
							"txid":              big.NewInt(0),
							"postings": []any{
								map[string]any{
									"source":      "world",
									"destination": "bank",
									"asset":       "USD/2",
									"amount":      1000,
								},
							},
							"reference": "foo",
							"timestamp": now.Format(time.RFC3339),
							"ledger":    "foo",
						},
					},
				},
			})
		})
		It("should be ingested properly", func() {
			Eventually(func(g Gomega) bool {
				response, err := Client().Search.V1.Search(
					TestContext(),
					shared.Query{
						Target: ptr("TRANSACTION"),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data).To(HaveLen(1))
				g.Expect(res.Cursor.Data[0]).To(Equal(map[string]any{
					"metadata":  map[string]any{},
					"reference": "foo",
					"postings": []any{
						map[string]any{
							"source":      "world",
							"destination": "bank",
							"asset":       "USD/2",
							"amount":      float64(1000),
						},
					},
					"txid":      float64(0),
					"timestamp": now.Format(time.RFC3339),
					"ledger":    "foo",
				}))

				return true
			}).Should(BeTrue())
			Eventually(func(g Gomega) bool {
				response, err := Client().Search.V1.Search(
					TestContext(),
					shared.Query{
						Target: ptr("ACCOUNT"),
					},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))

				res := response.Response
				g.Expect(res.Cursor.Data).To(HaveLen(2))
				g.Expect(res.Cursor.Data).To(ContainElements(
					map[string]any{
						"address":  "world",
						"ledger":   "foo",
						"metadata": map[string]any{},
					},
					map[string]any{
						"address":  "bank",
						"ledger":   "foo",
						"metadata": map[string]any{},
					},
				))
				return true
			}).Should(BeTrue())
		})
		Then("adding metadata on a specific account", func() {
			BeforeEach(func() {
				PublishLedger(publish.EventMessage{
					Date:    time.Now(),
					App:     "ledger",
					Version: "v1",
					Type:    "SAVED_METADATA",
					Payload: map[string]any{
						"targetId":   "bank",
						"targetType": "ACCOUNT",
						"metadata": map[string]any{
							"custom": map[string]any{
								"foo": "bar",
							},
						},
						"ledger": "foo",
					},
				})
			})
			It("should be ok", func() {
				Eventually(func(g Gomega) bool {
					response, err := Client().Search.V1.Search(
						TestContext(),
						shared.Query{
							Target: ptr("ACCOUNT"),
						},
					)
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(response.StatusCode).To(Equal(200))

					res := response.Response
					g.Expect(res.Cursor.Data).To(HaveLen(2))
					g.Expect(res.Cursor.Data).To(ContainElements(
						map[string]any{
							"address": "bank",
							"ledger":  "foo",
							"metadata": map[string]any{
								"custom": `{"foo":"bar"}`,
							},
						},
					))
					return true
				}).Should(BeTrue())
			})
		})
	})
})

package suite

import (
	"fmt"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"
	"net/http"
	"sort"
	"time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = WithModules([]*Module{modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))

		_, err = Client().Ledger.V2.CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "another",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("listing logs", func() {
		var (
			timestamp1 = time.Date(2023, 4, 11, 10, 0, 0, 0, time.UTC)
			timestamp2 = time.Date(2023, 4, 12, 10, 0, 0, 0, time.UTC)

			m1 = map[string]string{
				"clientType": "silver",
			}
			m2 = map[string]string{
				"clientType": "gold",
			}
		)
		BeforeEach(func() {
			response, err := Client().Ledger.V2.CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "foo:foo",
							},
						},
						Timestamp: &timestamp1,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			_, err = Client().Ledger.V2.CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "foo:foo",
							},
						},
						Timestamp: &timestamp1,
					},
					Ledger: "another",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			response, err = Client().Ledger.V2.CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: m1,
						Postings: []shared.V2Posting{
							{
								Amount:      big.NewInt(200),
								Asset:       "USD",
								Source:      "world",
								Destination: "foo:bar",
							},
						},
						Timestamp: &timestamp2,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			addMetadataResponse, err := Client().Ledger.V2.AddMetadataToAccount(
				TestContext(),
				operations.V2AddMetadataToAccountRequest{
					RequestBody: m2,
					Address:     "foo:baz",
					Ledger:      "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(addMetadataResponse.StatusCode).To(Equal(204))
		})
		It("should be listed on api with ListLogs", func() {
			response, err := Client().Ledger.V2.ListLogs(
				TestContext(),
				operations.V2ListLogsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			logsCursorResponse := response.V2LogsCursorResponse
			Expect(logsCursorResponse.Cursor.Data).To(HaveLen(3))

			// Cannot check the date and the hash since they are changing at
			// every run
			Expect(logsCursorResponse.Cursor.Data[0].ID).To(Equal(big.NewInt(2)))
			Expect(logsCursorResponse.Cursor.Data[0].Type).To(Equal(shared.V2LogTypeSetMetadata))
			Expect(logsCursorResponse.Cursor.Data[0].Data).To(Equal(map[string]any{
				"targetType": "ACCOUNT",
				"metadata": map[string]any{
					"clientType": "gold",
				},
				"targetId": "foo:baz",
			}))

			Expect(logsCursorResponse.Cursor.Data[1].ID).To(Equal(big.NewInt((1))))
			Expect(logsCursorResponse.Cursor.Data[1].Type).To(Equal(shared.V2LogTypeNewTransaction))
			// Cannot check date and txid inside Data since they are changing at
			// every run
			Expect(logsCursorResponse.Cursor.Data[1].Data["accountMetadata"]).To(Equal(map[string]any{}))
			Expect(logsCursorResponse.Cursor.Data[1].Data["transaction"]).To(BeAssignableToTypeOf(map[string]any{}))
			transaction := logsCursorResponse.Cursor.Data[1].Data["transaction"].(map[string]any)
			Expect(transaction["metadata"]).To(Equal(map[string]any{
				"clientType": "silver",
			}))
			Expect(transaction["timestamp"]).To(Equal("2023-04-12T10:00:00Z"))
			Expect(transaction["postings"]).To(Equal([]any{
				map[string]any{
					"amount":      float64(200),
					"asset":       "USD",
					"source":      "world",
					"destination": "foo:bar",
				},
			}))

			Expect(logsCursorResponse.Cursor.Data[2].ID).To(Equal(big.NewInt((0))))
			Expect(logsCursorResponse.Cursor.Data[2].Type).To(Equal(shared.V2LogTypeNewTransaction))
			Expect(logsCursorResponse.Cursor.Data[2].Data["accountMetadata"]).To(Equal(map[string]any{}))
			Expect(logsCursorResponse.Cursor.Data[2].Data["transaction"]).To(BeAssignableToTypeOf(map[string]any{}))
			transaction = logsCursorResponse.Cursor.Data[2].Data["transaction"].(map[string]any)
			Expect(transaction["metadata"]).To(Equal(map[string]any{}))
			Expect(transaction["timestamp"]).To(Equal("2023-04-11T10:00:00Z"))
			Expect(transaction["postings"]).To(Equal([]any{
				map[string]any{
					"amount":      float64(100),
					"asset":       "USD",
					"source":      "world",
					"destination": "foo:foo",
				},
			}))
		})
	})

	type expectedLog struct {
		id       *big.Int
		typ      shared.V2LogType
		postings []any
	}

	var (
		compareLogs = func(log shared.V2Log, expected expectedLog) {
			Expect(log.ID).To(Equal(expected.id))
			Expect(log.Type).To(Equal(expected.typ))
			Expect(log.Data["accountMetadata"]).To(Equal(map[string]any{}))
			Expect(log.Data["transaction"]).To(BeAssignableToTypeOf(map[string]any{}))
			transaction := log.Data["transaction"].(map[string]any)
			Expect(transaction["metadata"]).To(Equal(map[string]any{}))
			Expect(transaction["postings"]).To(Equal(expected.postings))
		}
	)

	const (
		pageSize      = int64(10)
		accountCounts = 2 * pageSize
	)
	When("creating logs with transactions", func() {
		var (
			expectedLogs []expectedLog
		)
		BeforeEach(func() {
			for i := int64(0); i < accountCounts; i++ {
				now := time.Now().Round(time.Millisecond).UTC()

				response, err := Client().Ledger.V2.CreateTransaction(
					TestContext(),
					operations.V2CreateTransactionRequest{
						V2PostTransaction: shared.V2PostTransaction{
							Metadata: map[string]string{},
							Postings: []shared.V2Posting{
								{
									Amount:      big.NewInt(100),
									Asset:       "USD",
									Source:      "world",
									Destination: fmt.Sprintf("foo:%d", i),
								},
							},
							Timestamp: &now,
						},
						Ledger: "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				expectedLogs = append(expectedLogs, expectedLog{
					id:  big.NewInt(i),
					typ: shared.V2LogTypeNewTransaction,
					postings: []any{
						map[string]any{
							"amount":      float64(100),
							"asset":       "USD",
							"source":      "world",
							"destination": fmt.Sprintf("foo:%d", i),
						},
					},
				})
			}

			sort.Slice(expectedLogs, func(i, j int) bool {
				return expectedLogs[i].id.Cmp(expectedLogs[j].id) > 0
			})
		})
		AfterEach(func() {
			expectedLogs = nil
		})
		Then(fmt.Sprintf("listing accounts using page size of %d", pageSize), func() {
			var (
				rsp *shared.V2LogsCursorResponse
			)
			BeforeEach(func() {
				response, err := Client().Ledger.V2.ListLogs(
					TestContext(),
					operations.V2ListLogsRequest{
						Ledger:   "default",
						PageSize: ptr(pageSize),
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				rsp = response.V2LogsCursorResponse
				Expect(rsp.Cursor.HasMore).To(BeTrue())
				Expect(rsp.Cursor.Previous).To(BeNil())
				Expect(rsp.Cursor.Next).NotTo(BeNil())
			})
			It("should return the first page", func() {
				Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
				Expect(len(rsp.Cursor.Data)).To(Equal(len(expectedLogs[:pageSize])))
				for i := range rsp.Cursor.Data {
					compareLogs(rsp.Cursor.Data[i], expectedLogs[i])
				}
			})
			Then("following next cursor", func() {
				BeforeEach(func() {
					response, err := Client().Ledger.V2.ListLogs(
						TestContext(),
						operations.V2ListLogsRequest{
							Cursor: rsp.Cursor.Next,
							Ledger: "default",
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))

					rsp = response.V2LogsCursorResponse
				})
				It("should return next page", func() {
					Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
					Expect(len(rsp.Cursor.Data)).To(Equal(len(expectedLogs[pageSize : 2*pageSize])))
					for i := range rsp.Cursor.Data {
						compareLogs(rsp.Cursor.Data[i], expectedLogs[int64(i)+pageSize])
					}
					Expect(rsp.Cursor.Next).To(BeNil())
				})
				Then("following previous cursor", func() {
					BeforeEach(func() {
						response, err := Client().Ledger.V2.ListLogs(
							TestContext(),
							operations.V2ListLogsRequest{
								Cursor: rsp.Cursor.Previous,
								Ledger: "default",
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						rsp = response.V2LogsCursorResponse
					})
					It("should return first page", func() {
						Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
						Expect(len(rsp.Cursor.Data)).To(Equal(len(expectedLogs[:pageSize])))
						for i := range rsp.Cursor.Data {
							compareLogs(rsp.Cursor.Data[i], expectedLogs[i])
						}
						Expect(rsp.Cursor.Previous).To(BeNil())
					})
				})
			})
		})
	})
})

package suite

import (
	"fmt"
	"math/big"
	"net/http"
	"sort"
	"time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/tests/integration/internal/modules"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
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
	})
	When("counting and listing accounts", func() {
		var (
			metadata1 = map[string]string{
				"clientType": "gold",
			}

			metadata2 = map[string]string{
				"clientType": "silver",
			}

			timestamp = time.Now().Round(time.Second).UTC()
			bigInt, _ = big.NewInt(0).SetString("999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)
		)
		BeforeEach(func() {
			// Subscribe to nats subject
			response, err := Client().Ledger.V2.AddMetadataToAccount(
				TestContext(),
				operations.V2AddMetadataToAccountRequest{
					RequestBody: metadata1,
					Address:     "foo:foo",
					Ledger:      "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

			response, err = Client().Ledger.V2.AddMetadataToAccount(
				TestContext(),
				operations.V2AddMetadataToAccountRequest{
					RequestBody: metadata2,
					Address:     "foo:bar",
					Ledger:      "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

			createTransactionResponse, err := Client().Ledger.V2.CreateTransaction(
				TestContext(),
				operations.V2CreateTransactionRequest{
					V2PostTransaction: shared.V2PostTransaction{
						Metadata: map[string]string{},
						Postings: []shared.V2Posting{
							{
								Amount:      bigInt,
								Asset:       "USD",
								Source:      "world",
								Destination: "foo:foo",
							},
						},
						Timestamp: &timestamp,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(createTransactionResponse.StatusCode).To(Equal(200))
		})
		It("should return a "+string(shared.V2ErrorsEnumValidation)+" on invalid filter", func() {
			_, err := Client().Ledger.V2.ListAccounts(
				TestContext(),
				operations.V2ListAccountsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"invalid-key": 0,
						},
					},
				},
			)
			Expect(err).To(HaveOccurred())
			Expect(err.(*sdkerrors.V2ErrorResponse).ErrorCode).To(Equal(shared.V2ErrorsEnumValidation))
		})
		It("should be countable on api", func() {
			response, err := Client().Ledger.V2.CountAccounts(
				TestContext(),
				operations.V2CountAccountsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))
			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("3"))
		})
		It("should be listed on api", func() {
			response, err := Client().Ledger.V2.ListAccounts(
				TestContext(),
				operations.V2ListAccountsRequest{
					Ledger: "default",
					Expand: pointer.For("volumes"),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse := response.V2AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(3))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.V2Account{
				Address:  "foo:bar",
				Metadata: metadata2,
			}))
			Expect(accountsCursorResponse.Cursor.Data[1]).To(Equal(shared.V2Account{
				Address:  "foo:foo",
				Metadata: metadata1,
				Volumes: map[string]shared.V2Volume{
					"USD": {
						Input:   bigInt,
						Output:  big.NewInt(0),
						Balance: bigInt,
					},
				},
			}))
			Expect(accountsCursorResponse.Cursor.Data[2]).To(Equal(shared.V2Account{
				Address:  "world",
				Metadata: metadata.Metadata{},
				Volumes: map[string]shared.V2Volume{
					"USD": {
						Output:  bigInt,
						Input:   big.NewInt(0),
						Balance: big.NewInt(0).Neg(bigInt),
					},
				},
			}))
		})
		It("should be listed on api using address filters", func() {
			response, err := Client().Ledger.V2.ListAccounts(
				TestContext(),
				operations.V2ListAccountsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"address": "foo:",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse := response.V2AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(2))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.V2Account{
				Address:  "foo:bar",
				Metadata: metadata2,
			}))
			Expect(accountsCursorResponse.Cursor.Data[1]).To(Equal(shared.V2Account{
				Address:  "foo:foo",
				Metadata: metadata1,
			}))

			response, err = Client().Ledger.V2.ListAccounts(
				TestContext(),
				operations.V2ListAccountsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"address": ":foo",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse = response.V2AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.V2Account{
				Address:  "foo:foo",
				Metadata: metadata1,
			}))
		})
		It("should be listed on api using metadata filters", func() {
			response, err := Client().Ledger.V2.ListAccounts(
				TestContext(),
				operations.V2ListAccountsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$match": map[string]any{
							"metadata[clientType]": "gold",
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse := response.V2AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(1))
			Expect(accountsCursorResponse.Cursor.Data[0]).To(Equal(shared.V2Account{
				Address:  "foo:foo",
				Metadata: metadata1,
			}))
		})
		It("should be listable on api using $not filter", func() {
			response, err := Client().Ledger.V2.ListAccounts(
				TestContext(),
				operations.V2ListAccountsRequest{
					Ledger: "default",
					RequestBody: map[string]interface{}{
						"$not": map[string]any{
							"$match": map[string]any{
								"address": "world",
							},
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			accountsCursorResponse := response.V2AccountsCursorResponse
			Expect(accountsCursorResponse.Cursor.Data).To(HaveLen(2))
		})
	})

	When("counting and listing accounts empty", func() {
		It("should be countable on api even if empty", func() {
			response, err := Client().Ledger.V2.CountAccounts(
				TestContext(),
				operations.V2CountAccountsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(204))

			Expect(response.Headers["Count"]).Should(HaveLen(1))
			Expect(response.Headers["Count"][0]).Should(Equal("0"))
		})
		It("should be listed on api even if empty", func() {
			response, err := Client().Ledger.V2.ListAccounts(
				TestContext(),
				operations.V2ListAccountsRequest{
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			Expect(response.V2AccountsCursorResponse.Cursor.Data).To(HaveLen(0))
		})
	})

	const (
		pageSize      = int64(10)
		accountCounts = 2 * pageSize
	)
	When("creating accounts", func() {
		var (
			accounts []shared.V2Account
		)
		BeforeEach(func() {
			for i := 0; i < int(accountCounts); i++ {
				m := map[string]string{
					"id": fmt.Sprintf("%d", i),
				}

				response, err := Client().Ledger.V2.AddMetadataToAccount(
					TestContext(),
					operations.V2AddMetadataToAccountRequest{
						RequestBody: m,
						Address:     fmt.Sprintf("foo:%d", i),
						Ledger:      "default",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(204))

				accounts = append(accounts, shared.V2Account{
					Address:  fmt.Sprintf("foo:%d", i),
					Metadata: m,
				})

				sort.Slice(accounts, func(i, j int) bool {
					return accounts[i].Address < accounts[j].Address
				})
			}
		})
		AfterEach(func() {
			accounts = nil
		})
		Then(fmt.Sprintf("listing accounts using page size of %d", pageSize), func() {
			var (
				rsp *shared.V2AccountsCursorResponse
			)
			BeforeEach(func() {
				response, err := Client().Ledger.V2.ListAccounts(
					TestContext(),
					operations.V2ListAccountsRequest{
						Ledger:   "default",
						PageSize: ptr(pageSize),
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				rsp = response.V2AccountsCursorResponse
				Expect(rsp.Cursor.HasMore).To(BeTrue())
				Expect(rsp.Cursor.Previous).To(BeNil())
				Expect(rsp.Cursor.Next).NotTo(BeNil())
			})
			It("should return the first page", func() {
				Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
				Expect(rsp.Cursor.Data).To(Equal(accounts[:pageSize]))
			})
			Then("following next cursor", func() {
				BeforeEach(func() {
					response, err := Client().Ledger.V2.ListAccounts(
						TestContext(),
						operations.V2ListAccountsRequest{
							Cursor: rsp.Cursor.Next,
							Ledger: "default",
						},
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(response.StatusCode).To(Equal(200))

					rsp = response.V2AccountsCursorResponse
				})
				It("should return next page", func() {
					Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
					Expect(rsp.Cursor.Data).To(Equal(accounts[pageSize : 2*pageSize]))
					Expect(rsp.Cursor.Next).To(BeNil())
				})
				Then("following previous cursor", func() {
					BeforeEach(func() {
						response, err := Client().Ledger.V2.ListAccounts(
							TestContext(),
							operations.V2ListAccountsRequest{
								Ledger: "default",
								Cursor: rsp.Cursor.Previous,
							},
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(response.StatusCode).To(Equal(200))

						rsp = response.V2AccountsCursorResponse
					})
					It("should return first page", func() {
						Expect(rsp.Cursor.PageSize).To(Equal(pageSize))
						Expect(rsp.Cursor.Data).To(Equal(accounts[:pageSize]))
						Expect(rsp.Cursor.Previous).To(BeNil())
					})
				})
			})
		})
	})

	When("Inserting one transaction in past and one in the future", func() {
		now := time.Now()
		BeforeEach(func() {
			_, err := Client().Ledger.V2.CreateTransaction(TestContext(), operations.V2CreateTransactionRequest{
				V2PostTransaction: shared.V2PostTransaction{
					Postings: []shared.V2Posting{{
						Amount:      big.NewInt(100),
						Asset:       "USD",
						Destination: "foo",
						Source:      "world",
					}},
					Timestamp: pointer.For(now.Add(-12 * time.Hour)),
					Metadata:  map[string]string{},
				},
				Ledger: "default",
			})
			Expect(err).To(Succeed())

			_, err = Client().Ledger.V2.CreateTransaction(TestContext(), operations.V2CreateTransactionRequest{
				V2PostTransaction: shared.V2PostTransaction{
					Postings: []shared.V2Posting{{
						Amount:      big.NewInt(100),
						Asset:       "USD",
						Destination: "foo",
						Source:      "world",
					}},
					Timestamp: pointer.For(now.Add(12 * time.Hour)),
					Metadata:  map[string]string{},
				},
				Ledger: "default",
			})
			Expect(err).To(Succeed())
		})
		When("getting account in the present", func() {
			It("should ignore future transaction on effective volumes", func() {
				accountResponse, err := Client().Ledger.V2.GetAccount(TestContext(), operations.V2GetAccountRequest{
					Address: "foo",
					Expand:  pointer.For("effectiveVolumes"),
					Ledger:  "default",
					Pit:     pointer.For(time.Now().Add(time.Minute)),
				})
				Expect(err).To(Succeed())
				Expect(accountResponse.V2AccountResponse.Data.EffectiveVolumes["USD"].Balance).To(Equal(big.NewInt(100)))
			})
		})
	})
})

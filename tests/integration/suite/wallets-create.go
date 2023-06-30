package suite

import (
	"fmt"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	const countWallets = 3
	When(fmt.Sprintf("creating %d wallets", countWallets), func() {
		JustBeforeEach(func() {
			for i := 0; i < countWallets; i++ {
				response, err := Client().Wallets.CreateWallet(
					TestContext(),
					shared.CreateWalletRequest{
						Metadata: map[string]string{
							"wallets_number": fmt.Sprint(i),
						},
						Name: uuid.NewString(),
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(201))
			}
		})
		Then("listing them", func() {
			var (
				request  operations.ListWalletsRequest
				response *operations.ListWalletsResponse
				err      error
			)
			JustBeforeEach(func() {
				Eventually(func(g Gomega) bool {
					response, err = Client().Wallets.ListWallets(TestContext(), request)
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(response.StatusCode).To(Equal(200))
					return true
				}).Should(BeTrue())
			})
			It(fmt.Sprintf("should return %d items", countWallets), func() {
				Expect(response.ListWalletsResponse.Cursor.Data).To(HaveLen(countWallets))
			})
			Context("using a metadata filter", func() {
				BeforeEach(func() {
					request.Metadata = map[string]string{
						"wallets_number": "0",
					}
				})
				It("should return only one item", func() {
					Expect(response.ListWalletsResponse.Cursor.Data).To(HaveLen(1))
				})
			})
		})
	})
})

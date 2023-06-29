package suite

import (
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("creating a new wallet", func() {
		BeforeEach(func() {
			response, err := Client().Wallets.CreateWallet(
				TestContext(),
				shared.CreateWalletRequest{
					Metadata: map[string]string{},
					Name:     uuid.NewString(),
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))
		})
		It("should be ok", func() {
			Eventually(func(g Gomega) []shared.Wallet {
				response, err := Client().Wallets.ListWallets(
					TestContext(),
					operations.ListWalletsRequest{},
				)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(response.StatusCode).To(Equal(200))
				return response.ListWalletsResponse.Cursor.Data
			}).Should(HaveLen(1)) // TODO: Check other fields
		})
	})
})

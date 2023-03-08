package suite

import (
	"github.com/formancehq/formance-sdk-go"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("creating a new wallet", func() {
		BeforeEach(func() {
			_, _, err := Client().WalletsApi.
				CreateWallet(TestContext()).
				CreateWalletRequest(formance.CreateWalletRequest{
					Name: uuid.NewString(),
				}).
				Execute()
			Expect(err).To(BeNil())
		})
		It("should be ok", func() {
			Eventually(func(g Gomega) []formance.Wallet {
				res, _, err := Client().WalletsApi.
					ListWallets(TestContext()).
					Execute()
				g.Expect(err).To(BeNil())
				return res.Cursor.Data
			}).Should(HaveLen(1)) // TODO: Check other fields
		})
	})
})

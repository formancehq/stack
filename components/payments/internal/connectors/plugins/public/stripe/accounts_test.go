package stripe_test

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/connectors/plugins/public/stripe"
	"github.com/formancehq/payments/internal/connectors/plugins/public/stripe/client"
	"github.com/formancehq/payments/internal/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	stripesdk "github.com/stripe/stripe-go/v79"
	gomock "go.uber.org/mock/gomock"
)

var _ = Describe("Stripe Plugin Accounts", func() {
	var (
		plg *stripe.Plugin
	)

	BeforeEach(func() {
		plg = &stripe.Plugin{}
	})

	Context("fetch next accounts", func() {
		var (
			m *client.MockClient

			pageSize       int
			sampleAccounts []*stripesdk.Account
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			m = client.NewMockClient(ctrl)
			plg.SetClient(m)
			pageSize = 20

			sampleAccounts = make([]*stripesdk.Account, 0)
			for i := 0; i < pageSize-1; i++ {
				sampleAccounts = append(sampleAccounts, &stripesdk.Account{
					ID: fmt.Sprintf("some-reference-%d", i),
				})
			}

		})
		It("fetches next accounts", func(ctx SpecContext) {
			req := models.FetchNextAccountsRequest{
				State:    json.RawMessage(`{}`),
				PageSize: pageSize,
			}
			// pageSize passed to client is less when we generate a root account
			m.EXPECT().GetAccounts(ctx, gomock.Any(), int64(pageSize-1)).Return(
				sampleAccounts,
				true,
				nil,
			)
			res, err := plg.FetchNextAccounts(ctx, req)
			Expect(err).To(BeNil())
			Expect(res.HasMore).To(BeTrue())
			Expect(res.Accounts).To(HaveLen(req.PageSize))
			Expect(res.Accounts[0].Reference).To(Equal("root"))

			var state stripe.AccountsState

			err = json.Unmarshal(res.NewState, &state)
			Expect(err).To(BeNil())
			Expect(state.LastID).To(Equal(res.Accounts[len(res.Accounts)-1].Reference))
		})
	})
})

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

var _ = Describe("Stripe Plugin ExternalAccounts", func() {
	var (
		plg *stripe.Plugin
	)

	BeforeEach(func() {
		plg = &stripe.Plugin{}
	})

	Context("fetch next ExternalAccounts", func() {
		var (
			m *client.MockClient

			sampleExternalAccounts []*stripesdk.BankAccount
			accRef                 string
			created                int64
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			m = client.NewMockClient(ctrl)
			plg.SetClient(m)

			accRef = "baseAcc"
			created = 1483565364
			sampleExternalAccounts = make([]*stripesdk.BankAccount, 0)
			for i := 0; i < int(stripe.PageLimit); i++ {
				sampleExternalAccounts = append(sampleExternalAccounts, &stripesdk.BankAccount{
					ID:      fmt.Sprintf("some-reference-%d", i),
					Account: &stripesdk.Account{Created: created},
				})
			}

		})
		It("fetches next ExternalAccounts", func(ctx SpecContext) {
			req := models.FetchNextExternalAccountsRequest{
				FromPayload: json.RawMessage(fmt.Sprintf(`{"reference": "%s"}`, accRef)),
				State:       json.RawMessage(`{}`),
			}
			m.EXPECT().GetExternalAccounts(ctx, &accRef, gomock.Any(), stripe.PageLimit).Return(
				sampleExternalAccounts,
				true,
				nil,
			)
			res, err := plg.FetchNextExternalAccounts(ctx, req)
			Expect(err).To(BeNil())
			Expect(res.HasMore).To(BeTrue())
			Expect(res.ExternalAccounts).To(HaveLen(int(stripe.PageLimit)))

			var state stripe.AccountsState

			err = json.Unmarshal(res.NewState, &state)
			Expect(err).To(BeNil())
			Expect(state.LastID).To(Equal(res.ExternalAccounts[len(res.ExternalAccounts)-1].Reference))
		})
	})
})

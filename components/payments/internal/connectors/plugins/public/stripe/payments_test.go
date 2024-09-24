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

var _ = Describe("Stripe Plugin Payments", func() {
	var (
		plg *stripe.Plugin
	)

	BeforeEach(func() {
		plg = &stripe.Plugin{}
	})

	Context("fetch next Payments", func() {
		var (
			m *client.MockClient

			samplePayments []*stripesdk.BalanceTransaction
			accRef         string
		)

		BeforeEach(func() {
			ctrl := gomock.NewController(GinkgoT())
			m = client.NewMockClient(ctrl)
			plg.SetClient(m)

			accRef = "baseAcc"
			samplePayments = make([]*stripesdk.BalanceTransaction, 0)
			for i := 0; i < int(stripe.PageLimit); i++ {
				samplePayments = append(samplePayments, &stripesdk.BalanceTransaction{
					ID: fmt.Sprintf("some-reference-%d", i),
				})
			}

		})
		It("fails when payments missing source are received", func(ctx SpecContext) {
			req := models.FetchNextPaymentsRequest{
				FromPayload: json.RawMessage(fmt.Sprintf(`{"reference": "%s"}`, accRef)),
				State:       json.RawMessage(`{}`),
			}
			m.EXPECT().GetPayments(ctx, &accRef, gomock.Any(), stripe.PageLimit).Return(
				samplePayments,
				true,
				nil,
			)
			res, err := plg.FetchNextPayments(ctx, req)
			Expect(err).To(MatchError(ContainSubstring(stripe.ErrInvalidPaymentSource.Error())))
			Expect(res.HasMore).To(BeFalse())
		})
	})
})

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
			samplePayments = []*stripesdk.BalanceTransaction{
				{
					ID:   "charge",
					Type: stripesdk.BalanceTransactionTypeCharge,
					Source: &stripesdk.BalanceTransactionSource{
						Charge: &stripesdk.Charge{
							Currency:             stripesdk.CurrencyBIF,
							PaymentMethodDetails: &stripesdk.ChargePaymentMethodDetails{Card: &stripesdk.ChargePaymentMethodDetailsCard{Brand: stripesdk.PaymentMethodCardBrandVisa}},
						},
					},
				},
				{
					ID:   "refund",
					Type: stripesdk.BalanceTransactionTypeRefund,
					Source: &stripesdk.BalanceTransactionSource{
						Refund: &stripesdk.Refund{
							Currency: stripesdk.CurrencyEUR,
							Charge: &stripesdk.Charge{
								Currency:             stripesdk.CurrencyGBP,
								BalanceTransaction:   &stripesdk.BalanceTransaction{ID: "refund_original"},
								PaymentMethodDetails: &stripesdk.ChargePaymentMethodDetails{Card: &stripesdk.ChargePaymentMethodDetailsCard{Brand: stripesdk.PaymentMethodCardBrandJCB}},
							},
						},
					},
				},
				{
					ID:   "refund_failure",
					Type: stripesdk.BalanceTransactionTypeRefundFailure,
					Source: &stripesdk.BalanceTransactionSource{
						Refund: &stripesdk.Refund{
							Currency: stripesdk.CurrencyGEL,
							Charge: &stripesdk.Charge{
								Currency:             stripesdk.CurrencyGIP,
								BalanceTransaction:   &stripesdk.BalanceTransaction{ID: "refund_failure_original"},
								PaymentMethodDetails: &stripesdk.ChargePaymentMethodDetails{Card: &stripesdk.ChargePaymentMethodDetailsCard{Brand: stripesdk.PaymentMethodCardBrandAmex}},
							},
						},
					},
				},
				{
					ID:   "payment",
					Type: stripesdk.BalanceTransactionTypePayment,
					Source: &stripesdk.BalanceTransactionSource{
						Charge: &stripesdk.Charge{
							Currency: stripesdk.CurrencyHKD,
						},
					},
				},
				{
					ID:   "payment_refund",
					Type: stripesdk.BalanceTransactionTypePaymentRefund,
					Source: &stripesdk.BalanceTransactionSource{
						Refund: &stripesdk.Refund{
							Currency: stripesdk.CurrencyEUR,
							Charge: &stripesdk.Charge{
								Currency:           stripesdk.CurrencyGBP,
								BalanceTransaction: &stripesdk.BalanceTransaction{ID: "payment_refund_original"},
							},
						},
					},
				},
				{
					ID:   "payment_refund_failure",
					Type: stripesdk.BalanceTransactionTypePaymentFailureRefund,
					Source: &stripesdk.BalanceTransactionSource{
						Refund: &stripesdk.Refund{
							Currency: stripesdk.CurrencyGEL,
							Charge: &stripesdk.Charge{
								Currency:           stripesdk.CurrencyGIP,
								BalanceTransaction: &stripesdk.BalanceTransaction{ID: "payment_refund_failure_original"},
							},
						},
					},
				},
				{
					ID:   "payout",
					Type: stripesdk.BalanceTransactionTypePayout,
					Source: &stripesdk.BalanceTransactionSource{
						Payout: &stripesdk.Payout{
							Currency:    stripesdk.CurrencyLKR,
							Status:      stripesdk.PayoutStatusPaid,
							Destination: &stripesdk.PayoutDestination{Card: &stripesdk.Card{Brand: stripesdk.CardBrandJCB}},
						},
					},
				},
				{
					ID:   "payout_failure",
					Type: stripesdk.BalanceTransactionTypePayoutFailure,
					Source: &stripesdk.BalanceTransactionSource{
						Payout: &stripesdk.Payout{
							Currency:           stripesdk.CurrencyMYR,
							BalanceTransaction: &stripesdk.BalanceTransaction{ID: "payout_failure_original"},
							Destination:        &stripesdk.PayoutDestination{Card: &stripesdk.Card{Brand: stripesdk.CardBrandUnionPay}},
						},
					},
				},
				{
					ID:   "transfer",
					Type: stripesdk.BalanceTransactionTypeTransfer,
					Source: &stripesdk.BalanceTransactionSource{
						Transfer: &stripesdk.Transfer{
							Currency: stripesdk.CurrencyCLP,
						},
					},
				},
				{
					ID:   "transfer_refund",
					Type: stripesdk.BalanceTransactionTypeTransferRefund,
					Source: &stripesdk.BalanceTransactionSource{
						Transfer: &stripesdk.Transfer{
							Currency:           stripesdk.CurrencyCOP,
							BalanceTransaction: &stripesdk.BalanceTransaction{ID: "transfer_refund_original"},
						},
					},
				},
				{
					ID:   "transfer_failed",
					Type: stripesdk.BalanceTransactionTypeTransferFailure,
					Source: &stripesdk.BalanceTransactionSource{
						Transfer: &stripesdk.Transfer{
							Currency:           stripesdk.CurrencyCAD,
							BalanceTransaction: &stripesdk.BalanceTransaction{ID: "transfer_failed_original"},
						},
					},
				},
				{
					ID:   "adjustment_with_dispute",
					Type: stripesdk.BalanceTransactionTypeAdjustment,
					Source: &stripesdk.BalanceTransactionSource{
						Dispute: &stripesdk.Dispute{
							Charge: &stripesdk.Charge{
								Currency:           stripesdk.CurrencyDZD,
								BalanceTransaction: &stripesdk.BalanceTransaction{ID: "adjustment_with_dispute_original"},
							},
						},
					},
				},
				{
					ID:   "skipped", // unsupported types are skipped
					Type: stripesdk.BalanceTransactionTypeStripeFee,
					Source: &stripesdk.BalanceTransactionSource{
						Charge: &stripesdk.Charge{
							Currency: stripesdk.CurrencyJPY,
						},
					},
				},
			}

		})
		It("fails when payments missing source are received", func(ctx SpecContext) {
			req := models.FetchNextPaymentsRequest{
				FromPayload: json.RawMessage(fmt.Sprintf(`{"reference": "%s"}`, accRef)),
				State:       json.RawMessage(`{}`),
			}
			p := []*stripesdk.BalanceTransaction{
				{
					ID:   "someid",
					Type: stripesdk.BalanceTransactionTypeAdjustment,
				},
			}
			m.EXPECT().GetPayments(ctx, &accRef, gomock.Any(), stripe.PageLimit).Return(
				p,
				true,
				nil,
			)
			res, err := plg.FetchNextPayments(ctx, req)
			Expect(err).To(MatchError(ContainSubstring(stripe.ErrInvalidPaymentSource.Error())))
			Expect(res.HasMore).To(BeFalse())
		})

		It("fails when payments contain unsupported currencies", func(ctx SpecContext) {
			req := models.FetchNextPaymentsRequest{
				FromPayload: json.RawMessage(fmt.Sprintf(`{"reference": "%s"}`, accRef)),
				State:       json.RawMessage(`{}`),
			}
			p := []*stripesdk.BalanceTransaction{
				{
					ID:   "someid",
					Type: stripesdk.BalanceTransactionTypeCharge,
					Source: &stripesdk.BalanceTransactionSource{
						Charge: &stripesdk.Charge{
							Currency: stripesdk.CurrencyEEK,
						},
					},
				},
			}
			m.EXPECT().GetPayments(ctx, &accRef, gomock.Any(), stripe.PageLimit).Return(
				p,
				true,
				nil,
			)
			res, err := plg.FetchNextPayments(ctx, req)
			Expect(err).To(MatchError(ContainSubstring(stripe.ErrUnsupportedCurrency.Error())))
			Expect(res.HasMore).To(BeFalse())
		})

		It("fetches payments", func(ctx SpecContext) {
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
			Expect(err).To(BeNil())
			Expect(res.Payments).To(HaveLen(len(samplePayments) - 1))
			Expect(res.HasMore).To(BeTrue())

			// Charges
			Expect(res.Payments[0].Reference).To(Equal(samplePayments[0].ID))
			Expect(res.Payments[0].Type).To(Equal(models.PAYMENT_TYPE_PAYIN))
			Expect(res.Payments[0].Status).To(Equal(models.PAYMENT_STATUS_SUCCEEDED))
			Expect(res.Payments[1].Reference).To(Equal(samplePayments[1].Source.Refund.Charge.BalanceTransaction.ID))
			Expect(res.Payments[1].Type).To(Equal(models.PAYMENT_TYPE_PAYIN))
			Expect(res.Payments[1].Status).To(Equal(models.PAYMENT_STATUS_REFUNDED))
			Expect(res.Payments[2].Reference).To(Equal(samplePayments[2].Source.Refund.Charge.BalanceTransaction.ID))
			Expect(res.Payments[2].Type).To(Equal(models.PAYMENT_TYPE_PAYIN))
			Expect(res.Payments[2].Status).To(Equal(models.PAYMENT_STATUS_REFUNDED_FAILURE))
			// Payments
			Expect(res.Payments[3].Reference).To(Equal(samplePayments[3].ID))
			Expect(res.Payments[3].Type).To(Equal(models.PAYMENT_TYPE_PAYIN))
			Expect(res.Payments[3].Status).To(Equal(models.PAYMENT_STATUS_SUCCEEDED))
			Expect(res.Payments[4].Reference).To(Equal(samplePayments[4].Source.Refund.Charge.BalanceTransaction.ID))
			Expect(res.Payments[4].Type).To(Equal(models.PAYMENT_TYPE_PAYIN))
			Expect(res.Payments[4].Status).To(Equal(models.PAYMENT_STATUS_REFUNDED))
			Expect(res.Payments[5].Reference).To(Equal(samplePayments[5].Source.Refund.Charge.BalanceTransaction.ID))
			Expect(res.Payments[5].Type).To(Equal(models.PAYMENT_TYPE_PAYIN))
			Expect(res.Payments[5].Status).To(Equal(models.PAYMENT_STATUS_REFUNDED_FAILURE))
			// Payouts
			Expect(res.Payments[6].Reference).To(Equal(samplePayments[6].ID))
			Expect(res.Payments[6].Type).To(Equal(models.PAYMENT_TYPE_PAYOUT))
			Expect(res.Payments[6].Status).To(Equal(models.PAYMENT_STATUS_SUCCEEDED))
			Expect(res.Payments[7].Reference).To(Equal(samplePayments[7].Source.Payout.BalanceTransaction.ID))
			Expect(res.Payments[7].Type).To(Equal(models.PAYMENT_TYPE_PAYOUT))
			Expect(res.Payments[7].Status).To(Equal(models.PAYMENT_STATUS_FAILED))
			// Transfers
			Expect(res.Payments[8].Reference).To(Equal(samplePayments[8].ID))
			Expect(res.Payments[8].Type).To(Equal(models.PAYMENT_TYPE_TRANSFER))
			Expect(res.Payments[8].Status).To(Equal(models.PAYMENT_STATUS_SUCCEEDED))
			Expect(res.Payments[9].Reference).To(Equal(samplePayments[9].Source.Transfer.BalanceTransaction.ID))
			Expect(res.Payments[9].Type).To(Equal(models.PAYMENT_TYPE_TRANSFER))
			Expect(res.Payments[9].Status).To(Equal(models.PAYMENT_STATUS_REFUNDED))
			Expect(res.Payments[10].Reference).To(Equal(samplePayments[10].Source.Transfer.BalanceTransaction.ID))
			Expect(res.Payments[10].Type).To(Equal(models.PAYMENT_TYPE_TRANSFER))
			Expect(res.Payments[10].Status).To(Equal(models.PAYMENT_STATUS_FAILED))
			// Adjustments
			Expect(res.Payments[11].Reference).To(Equal(samplePayments[11].Source.Dispute.Charge.BalanceTransaction.ID))
			Expect(res.Payments[11].Type).To(Equal(models.PAYMENT_TYPE_PAYIN))
			Expect(res.Payments[11].Status).To(Equal(models.PAYMENT_STATUS_DISPUTE))

			var state stripe.PaymentState

			err = json.Unmarshal(res.NewState, &state)
			Expect(err).To(BeNil())
			Expect(state.LastID).To(Equal(samplePayments[len(samplePayments)-1].ID))
		})
	})
})

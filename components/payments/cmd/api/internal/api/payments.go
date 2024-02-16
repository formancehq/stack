package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/gorilla/mux"
)

type paymentResponse struct {
	ID                   string                   `json:"id"`
	Reference            string                   `json:"reference"`
	SourceAccountID      string                   `json:"sourceAccountID"`
	DestinationAccountID string                   `json:"destinationAccountID"`
	Type                 string                   `json:"type"`
	Provider             models.ConnectorProvider `json:"provider"`
	ConnectorID          string                   `json:"connectorID"`
	Status               models.PaymentStatus     `json:"status"`
	Amount               *big.Int                 `json:"amount"`
	InitialAmount        *big.Int                 `json:"initialAmount"`
	Scheme               models.PaymentScheme     `json:"scheme"`
	Asset                string                   `json:"asset"`
	CreatedAt            time.Time                `json:"createdAt"`
	Raw                  interface{}              `json:"raw"`
	Adjustments          []paymentAdjustment      `json:"adjustments"`
	Metadata             map[string]string        `json:"metadata"`
}

type paymentAdjustment struct {
	Reference string               `json:"reference" bson:"reference"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	Status    models.PaymentStatus `json:"status" bson:"status"`
	Amount    *big.Int             `json:"amount" bson:"amount"`
	Raw       interface{}          `json:"raw" bson:"raw"`
}

func createPaymentHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req service.CreatePaymentRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		if err := req.Validate(); err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		payment, err := b.GetService().CreatePayment(r.Context(), &req)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := paymentResponse{
			ID:            payment.ID.String(),
			Reference:     payment.Reference,
			Type:          payment.Type.String(),
			ConnectorID:   payment.ConnectorID.String(),
			Provider:      payment.ConnectorID.Provider,
			Status:        payment.Status,
			Amount:        payment.Amount,
			InitialAmount: payment.InitialAmount,
			Scheme:        payment.Scheme,
			Asset:         payment.Asset.String(),
			CreatedAt:     payment.CreatedAt,
			Raw:           payment.RawData,
			Adjustments:   make([]paymentAdjustment, len(payment.Adjustments)),
		}

		if payment.SourceAccountID != nil {
			data.SourceAccountID = payment.SourceAccountID.String()
		}

		if payment.DestinationAccountID != nil {
			data.DestinationAccountID = payment.DestinationAccountID.String()
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[paymentResponse]{
			Data: &data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func listPaymentsHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		query, err := bunpaginate.Extract[storage.ListPaymentsQuery](r, func() (*storage.ListPaymentsQuery, error) {
			options, err := getPagination(r, storage.PaymentQuery{})
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListPaymentsQuery(*options)), nil
		})
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := b.GetService().ListPayments(r.Context(), *query)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		ret := cursor.Data
		data := make([]*paymentResponse, len(ret))

		for i := range ret {
			data[i] = &paymentResponse{
				ID:            ret[i].ID.String(),
				Reference:     ret[i].Reference,
				Type:          ret[i].Type.String(),
				ConnectorID:   ret[i].ConnectorID.String(),
				Provider:      ret[i].Connector.Provider,
				Status:        ret[i].Status,
				Amount:        ret[i].Amount,
				InitialAmount: ret[i].InitialAmount,
				Scheme:        ret[i].Scheme,
				Asset:         ret[i].Asset.String(),
				CreatedAt:     ret[i].CreatedAt,
				Raw:           ret[i].RawData,
				Adjustments:   make([]paymentAdjustment, len(ret[i].Adjustments)),
			}

			if ret[i].Connector != nil {
				data[i].Provider = ret[i].Connector.Provider
			}

			if ret[i].SourceAccountID != nil {
				data[i].SourceAccountID = ret[i].SourceAccountID.String()
			}

			if ret[i].DestinationAccountID != nil {
				data[i].DestinationAccountID = ret[i].DestinationAccountID.String()
			}

			for adjustmentIdx := range ret[i].Adjustments {
				data[i].Adjustments[adjustmentIdx] = paymentAdjustment{
					Reference: ret[i].Adjustments[adjustmentIdx].Reference,
					Status:    ret[i].Adjustments[adjustmentIdx].Status,
					Amount:    ret[i].Adjustments[adjustmentIdx].Amount,
					CreatedAt: ret[i].Adjustments[adjustmentIdx].CreatedAt,
					Raw:       ret[i].Adjustments[adjustmentIdx].RawData,
				}
			}

			if ret[i].Metadata != nil {
				data[i].Metadata = make(map[string]string)

				for metadataIDx := range ret[i].Metadata {
					data[i].Metadata[ret[i].Metadata[metadataIDx].Key] = ret[i].Metadata[metadataIDx].Value
				}
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*paymentResponse]{
			Cursor: &api.Cursor[*paymentResponse]{
				PageSize: cursor.PageSize,
				HasMore:  cursor.HasMore,
				Previous: cursor.Previous,
				Next:     cursor.Next,
				Data:     data,
			},
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func readPaymentHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		paymentID := mux.Vars(r)["paymentID"]

		payment, err := b.GetService().GetPayment(r.Context(), paymentID)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := paymentResponse{
			ID:            payment.ID.String(),
			Reference:     payment.Reference,
			Type:          payment.Type.String(),
			ConnectorID:   payment.ConnectorID.String(),
			Status:        payment.Status,
			Amount:        payment.Amount,
			InitialAmount: payment.InitialAmount,
			Scheme:        payment.Scheme,
			Asset:         payment.Asset.String(),
			CreatedAt:     payment.CreatedAt,
			Raw:           payment.RawData,
			Adjustments:   make([]paymentAdjustment, len(payment.Adjustments)),
		}

		if payment.SourceAccountID != nil {
			data.SourceAccountID = payment.SourceAccountID.String()
		}

		if payment.DestinationAccountID != nil {
			data.DestinationAccountID = payment.DestinationAccountID.String()
		}

		if payment.Connector != nil {
			data.Provider = payment.Connector.Provider
		}

		for i := range payment.Adjustments {
			data.Adjustments[i] = paymentAdjustment{
				Reference: payment.Adjustments[i].Reference,
				Status:    payment.Adjustments[i].Status,
				Amount:    payment.Adjustments[i].Amount,
				CreatedAt: payment.Adjustments[i].CreatedAt,
				Raw:       payment.Adjustments[i].RawData,
			}
		}

		if payment.Metadata != nil {
			data.Metadata = make(map[string]string)

			for metadataIDx := range payment.Metadata {
				data.Metadata[payment.Metadata[metadataIDx].Key] = payment.Metadata[metadataIDx].Value
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[paymentResponse]{
			Data: &data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

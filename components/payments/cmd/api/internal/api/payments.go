package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
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
	InitialAmount        *big.Int                 `json:"initialAmount"`
	Scheme               models.PaymentScheme     `json:"scheme"`
	Asset                string                   `json:"asset"`
	CreatedAt            time.Time                `json:"createdAt"`
	Raw                  interface{}              `json:"raw"`
	Adjustments          []paymentAdjustment      `json:"adjustments"`
	Metadata             map[string]string        `json:"metadata"`
}

type paymentAdjustment struct {
	Status   models.PaymentStatus `json:"status" bson:"status"`
	Amount   int64                `json:"amount" bson:"amount"`
	Date     time.Time            `json:"date" bson:"date"`
	Raw      interface{}          `json:"raw" bson:"raw"`
	Absolute bool                 `json:"absolute" bson:"absolute"`
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
			InitialAmount: payment.Amount,
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

		pagination, err := getPagination(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		ret, paginationDetails, err := b.GetService().ListPayments(r.Context(), pagination)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]*paymentResponse, len(ret))

		for i := range ret {
			data[i] = &paymentResponse{
				ID:            ret[i].ID.String(),
				Reference:     ret[i].Reference,
				Type:          ret[i].Type.String(),
				ConnectorID:   ret[i].ConnectorID.String(),
				Provider:      ret[i].Connector.Provider,
				Status:        ret[i].Status,
				InitialAmount: ret[i].Amount,
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
					Status:   ret[i].Adjustments[adjustmentIdx].Status,
					Amount:   ret[i].Adjustments[adjustmentIdx].Amount,
					Date:     ret[i].Adjustments[adjustmentIdx].CreatedAt,
					Raw:      ret[i].Adjustments[adjustmentIdx].RawData,
					Absolute: ret[i].Adjustments[adjustmentIdx].Absolute,
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
				PageSize: paginationDetails.PageSize,
				HasMore:  paginationDetails.HasMore,
				Previous: paginationDetails.PreviousPage,
				Next:     paginationDetails.NextPage,
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
			InitialAmount: payment.Amount,
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
				Status:   payment.Adjustments[i].Status,
				Amount:   payment.Adjustments[i].Amount,
				Date:     payment.Adjustments[i].CreatedAt,
				Raw:      payment.Adjustments[i].RawData,
				Absolute: payment.Adjustments[i].Absolute,
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

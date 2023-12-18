package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
)

type transferInitiationResponse struct {
	ID                   string    `json:"id"`
	Reference            string    `json:"reference"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	ScheduledAt          time.Time `json:"scheduledAt"`
	Description          string    `json:"description"`
	SourceAccountID      string    `json:"sourceAccountID"`
	DestinationAccountID string    `json:"destinationAccountID"`
	ConnectorID          string    `json:"connectorID"`
	Provider             string    `json:"provider"`
	Type                 string    `json:"type"`
	Amount               *big.Int  `json:"amount"`
	Asset                string    `json:"asset"`
	Status               string    `json:"status"`
	Error                string    `json:"error"`
}

type transferInitiationPaymentsResponse struct {
	PaymentID string    `json:"paymentID"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
	Error     string    `json:"error"`
}

type readTransferInitiationResponse struct {
	transferInitiationResponse
	RelatedPayments []*transferInitiationPaymentsResponse `json:"relatedPayments"`
}

func readTransferInitiationHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		transferID, err := models.TransferInitiationIDFromString(mux.Vars(r)["transferID"])
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		ret, err := b.GetService().ReadTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := &readTransferInitiationResponse{
			transferInitiationResponse: transferInitiationResponse{
				ID:                   ret.ID.String(),
				Reference:            ret.ID.Reference,
				CreatedAt:            ret.CreatedAt,
				UpdatedAt:            ret.UpdatedAt,
				ScheduledAt:          ret.ScheduledAt,
				Description:          ret.Description,
				SourceAccountID:      ret.SourceAccountID.String(),
				DestinationAccountID: ret.DestinationAccountID.String(),
				ConnectorID:          ret.ConnectorID.String(),
				Provider:             ret.Provider.String(),
				Type:                 ret.Type.String(),
				Amount:               ret.Amount,
				Asset:                ret.Asset.String(),
				Status:               ret.Status.String(),
				Error:                ret.Error,
			},
		}

		for _, payments := range ret.RelatedPayments {
			data.RelatedPayments = append(data.RelatedPayments, &transferInitiationPaymentsResponse{
				PaymentID: payments.PaymentID.String(),
				CreatedAt: payments.CreatedAt,
				Status:    payments.Status.String(),
				Error:     payments.Error,
			})
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[readTransferInitiationResponse]{
			Data: data,
		})
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func listTransferInitiationsHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		pagination, err := getPagination(r)
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		ret, paginationDetails, err := b.GetService().ListTransferInitiations(r.Context(), pagination)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]*transferInitiationResponse, len(ret))
		for i := range ret {
			data[i] = &transferInitiationResponse{
				ID:                   ret[i].ID.String(),
				Reference:            ret[i].ID.Reference,
				CreatedAt:            ret[i].CreatedAt,
				UpdatedAt:            ret[i].UpdatedAt,
				ScheduledAt:          ret[i].ScheduledAt,
				Description:          ret[i].Description,
				SourceAccountID:      ret[i].SourceAccountID.String(),
				DestinationAccountID: ret[i].DestinationAccountID.String(),
				Provider:             ret[i].Provider.String(),
				ConnectorID:          ret[i].ConnectorID.String(),
				Type:                 ret[i].Type.String(),
				Amount:               ret[i].Amount,
				Asset:                ret[i].Asset.String(),
				Status:               ret[i].Status.String(),
				Error:                ret[i].Error,
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*transferInitiationResponse]{
			Cursor: &api.Cursor[*transferInitiationResponse]{
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

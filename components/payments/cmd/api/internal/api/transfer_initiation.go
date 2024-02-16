package api

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/gorilla/mux"
)

type transferInitiationResponse struct {
	ID                   string            `json:"id"`
	Reference            string            `json:"reference"`
	CreatedAt            time.Time         `json:"createdAt"`
	ScheduledAt          time.Time         `json:"scheduledAt"`
	Description          string            `json:"description"`
	SourceAccountID      string            `json:"sourceAccountID"`
	DestinationAccountID string            `json:"destinationAccountID"`
	ConnectorID          string            `json:"connectorID"`
	Provider             string            `json:"provider"`
	Type                 string            `json:"type"`
	Amount               *big.Int          `json:"amount"`
	InitialAmount        *big.Int          `json:"initialAmount"`
	Asset                string            `json:"asset"`
	Status               string            `json:"status"`
	Error                string            `json:"error"`
	Metadata             map[string]string `json:"metadata"`
}

type transferInitiationPaymentsResponse struct {
	PaymentID string    `json:"paymentID"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
	Error     string    `json:"error"`
}

type transferInitiationAdjustmentsResponse struct {
	AdjustmentID string            `json:"adjustmentID"`
	CreatedAt    time.Time         `json:"createdAt"`
	Status       string            `json:"status"`
	Error        string            `json:"error"`
	Metadata     map[string]string `json:"metadata"`
}

type readTransferInitiationResponse struct {
	transferInitiationResponse
	RelatedPayments    []*transferInitiationPaymentsResponse    `json:"relatedPayments"`
	RelatedAdjustments []*transferInitiationAdjustmentsResponse `json:"relatedAdjustments"`
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
				ScheduledAt:          ret.ScheduledAt,
				Description:          ret.Description,
				SourceAccountID:      ret.SourceAccountID.String(),
				DestinationAccountID: ret.DestinationAccountID.String(),
				ConnectorID:          ret.ConnectorID.String(),
				Provider:             ret.Provider.String(),
				Type:                 ret.Type.String(),
				Amount:               ret.Amount,
				InitialAmount:        ret.InitialAmount,
				Asset:                ret.Asset.String(),
				Metadata:             ret.Metadata,
			},
		}

		if len(ret.RelatedAdjustments) > 0 {
			// Take the status and error from the last adjustment
			data.Status = ret.RelatedAdjustments[0].Status.String()
			data.Error = ret.RelatedAdjustments[0].Error
		}

		for _, adjustments := range ret.RelatedAdjustments {
			data.RelatedAdjustments = append(data.RelatedAdjustments, &transferInitiationAdjustmentsResponse{
				AdjustmentID: adjustments.ID.String(),
				CreatedAt:    adjustments.CreatedAt,
				Status:       adjustments.Status.String(),
				Error:        adjustments.Error,
				Metadata:     adjustments.Metadata,
			})
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

		query, err := bunpaginate.Extract[storage.ListTransferInitiationsQuery](r, func() (*storage.ListTransferInitiationsQuery, error) {
			options, err := getPagination(r, storage.TransferInitiationQuery{})
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListTransferInitiationsQuery(*options)), nil
		})
		if err != nil {
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := b.GetService().ListTransferInitiations(r.Context(), *query)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		ret := cursor.Data
		data := make([]*transferInitiationResponse, len(ret))
		for i := range ret {
			ret[i].SortRelatedAdjustments()
			data[i] = &transferInitiationResponse{
				ID:                   ret[i].ID.String(),
				Reference:            ret[i].ID.Reference,
				CreatedAt:            ret[i].CreatedAt,
				ScheduledAt:          ret[i].ScheduledAt,
				Description:          ret[i].Description,
				SourceAccountID:      ret[i].SourceAccountID.String(),
				DestinationAccountID: ret[i].DestinationAccountID.String(),
				Provider:             ret[i].Provider.String(),
				ConnectorID:          ret[i].ConnectorID.String(),
				Type:                 ret[i].Type.String(),
				Amount:               ret[i].Amount,
				InitialAmount:        ret[i].InitialAmount,
				Asset:                ret[i].Asset.String(),
				Metadata:             ret[i].Metadata,
			}

			if len(ret[i].RelatedAdjustments) > 0 {
				// Take the status and error from the last adjustment
				data[i].Status = ret[i].RelatedAdjustments[0].Status.String()
				data[i].Error = ret[i].RelatedAdjustments[0].Error
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*transferInitiationResponse]{
			Cursor: &api.Cursor[*transferInitiationResponse]{
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

package api

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/storage"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type listPaymentsRepository interface {
	ListPayments(ctx context.Context, pagination storage.Paginator) ([]*models.Payment, storage.PaginationDetails, error)
}

type paymentResponse struct {
	ID                   string                   `json:"id"`
	Reference            string                   `json:"reference"`
	SourceAccountID      string                   `json:"sourceAccountID"`
	DestinationAccountID string                   `json:"destinationAccountID"`
	Type                 string                   `json:"type"`
	ConnectorID          string                   `json:"connectorID"`
	Provider             models.ConnectorProvider `json:"provider"`
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

func listPaymentsHandler(repo listPaymentsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var sorter storage.Sorter

		if sortParams := r.URL.Query()["sort"]; sortParams != nil {
			for _, s := range sortParams {
				parts := strings.SplitN(s, ":", 2)

				var order storage.SortOrder

				if len(parts) > 1 {
					switch parts[1] {
					case "asc", "ASC":
						order = storage.SortOrderAsc
					case "dsc", "desc", "DSC", "DESC":
						order = storage.SortOrderDesc
					default:
						handleValidationError(w, r, errors.New("sort order not well specified, got "+parts[1]))

						return
					}
				}

				column := parts[0]

				sorter.Add(column, order)
			}
		}

		pageSize, err := pageSizeQueryParam(r)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		pagination, err := storage.Paginate(pageSize, r.URL.Query().Get("cursor"), sorter)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		ret, paginationDetails, err := repo.ListPayments(r.Context(), pagination)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		data := make([]*paymentResponse, len(ret))

		for i := range ret {
			data[i] = &paymentResponse{
				ID:            ret[i].ID.String(),
				Reference:     ret[i].Reference,
				Type:          ret[i].Type.String(),
				Provider:      ret[i].Connector.Provider,
				ConnectorID:   ret[i].Connector.ID.String(),
				Status:        ret[i].Status,
				InitialAmount: ret[i].Amount,
				Scheme:        ret[i].Scheme,
				Asset:         ret[i].Asset.String(),
				CreatedAt:     ret[i].CreatedAt,
				Raw:           ret[i].RawData,
				Adjustments:   make([]paymentAdjustment, len(ret[i].Adjustments)),
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
			handleServerError(w, r, err)

			return
		}
	}
}

type readPaymentRepository interface {
	GetPayment(ctx context.Context, id string) (*models.Payment, error)
}

func readPaymentHandler(repo readPaymentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		paymentID := mux.Vars(r)["paymentID"]

		payment, err := repo.GetPayment(r.Context(), paymentID)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		data := paymentResponse{
			ID:            payment.ID.String(),
			Reference:     payment.Reference,
			Type:          payment.Type.String(),
			ConnectorID:   payment.Connector.ID.String(),
			Provider:      payment.Connector.Provider,
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
			handleServerError(w, r, err)

			return
		}
	}
}

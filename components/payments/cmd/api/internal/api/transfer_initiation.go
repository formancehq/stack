package api

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type transferInitiationResponse struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	ScheduledAt          time.Time `json:"scheduledAt"`
	Description          string    `json:"description"`
	SourceAccountID      string    `json:"sourceAccountID"`
	DestinationAccountID string    `json:"destinationAccountID"`
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

type readTransferInitiationRepository interface {
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
}

func readTransferInitiationHandler(repo readTransferInitiationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		transferID, err := models.TransferInitiationIDFromString(mux.Vars(r)["transferID"])
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		ret, err := repo.ReadTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		type readTransferInitiationResponse struct {
			transferInitiationResponse
			RelatedPayments []*transferInitiationPaymentsResponse `json:"relatedPayments"`
		}

		data := &readTransferInitiationResponse{
			transferInitiationResponse: transferInitiationResponse{
				ID:                   ret.ID.String(),
				CreatedAt:            ret.CreatedAt,
				UpdatedAt:            ret.UpdatedAt,
				ScheduledAt:          ret.ScheduledAt,
				Description:          ret.Description,
				SourceAccountID:      ret.SourceAccountID.String(),
				DestinationAccountID: ret.DestinationAccountID.String(),
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
			handleServerError(w, r, err)

			return
		}
	}
}

type listTransferInitiationsRepository interface {
	ListTransferInitiations(ctx context.Context, pagination storage.PaginatorQuery) ([]*models.TransferInitiation, storage.PaginationDetails, error)
}

func listTransferInitiationsHandler(repo listTransferInitiationsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		qb, err := getQueryBuilder(r)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		var sorter storage.Sorter

		if sortParams := r.URL.Query()["sort"]; sortParams != nil {
			for _, s := range sortParams {
				parts := strings.SplitN(s, ":", 2)

				var order storage.SortOrder

				if len(parts) > 1 {
					//nolint:goconst // allow duplicate string
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

		pagination, err := storage.Paginate(pageSize, r.URL.Query().Get("cursor"), sorter, qb)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		ret, paginationDetails, err := repo.ListTransferInitiations(r.Context(), pagination)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		data := make([]*transferInitiationResponse, len(ret))
		for i := range ret {
			data[i] = &transferInitiationResponse{
				ID:                   ret[i].ID.String(),
				CreatedAt:            ret[i].CreatedAt,
				UpdatedAt:            ret[i].UpdatedAt,
				ScheduledAt:          ret[i].ScheduledAt,
				Description:          ret[i].Description,
				SourceAccountID:      ret[i].SourceAccountID.String(),
				DestinationAccountID: ret[i].DestinationAccountID.String(),
				Provider:             ret[i].Provider.String(),
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
			handleServerError(w, r, err)

			return
		}
	}
}

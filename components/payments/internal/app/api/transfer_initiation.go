package api

import (
	"context"
	"encoding/json"
	"fmt"
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

type paymentHandler func(ctx context.Context, transfer *models.TransferInitiation) error

type transferInitiationResponse struct {
	ID                   string    `json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
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

type createTransferInitiationRequest struct {
	Reference            string    `json:"reference"`
	CreatedAt            time.Time `json:"createdAt"`
	Description          string    `json:"description"`
	SourceAccountID      string    `json:"sourceAccountID"`
	DestinationAccountID string    `json:"destinationAccountID"`
	Provider             string    `json:"provider"`
	Type                 string    `json:"type"`
	Amount               *big.Int  `json:"amount"`
	Asset                string    `json:"asset"`
	Validated            bool      `json:"validated"`
}

func (r *createTransferInitiationRequest) Validate(repo createTransferInitiationRepository) error {
	if r.Reference == "" {
		return errors.New("uniqueRequestId is required")
	}

	if r.Description == "" {
		return errors.New("description is required")
	}

	if r.SourceAccountID != "" {
		_, err := models.AccountIDFromString(r.SourceAccountID)
		if err != nil {
			return err
		}
	}

	_, err := models.AccountIDFromString(r.DestinationAccountID)
	if err != nil {
		return err
	}

	_, err = models.TransferInitiationTypeFromString(r.Type)
	if err != nil {
		return err
	}

	_, err = models.ConnectorProviderFromString(r.Provider)
	if err != nil {
		return err
	}

	if r.Amount == nil {
		return errors.New("amount is required")
	}

	if r.Asset == "" {
		return errors.New("asset is required")
	}

	return nil
}

type createTransferInitiationRepository interface {
	CreateTransferInitiation(ctx context.Context, transferInitiation *models.TransferInitiation) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
	IsInstalled(ctx context.Context, provider models.ConnectorProvider) (bool, error)
}

func createTransferInitiationHandler(
	repo createTransferInitiationRepository,
	paymentHandlers map[models.ConnectorProvider]paymentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload := &createTransferInitiationRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			handleValidationError(w, r, err)

			return
		}

		if err := payload.Validate(repo); err != nil {
			handleValidationError(w, r, err)

			return
		}

		if payload.CreatedAt.IsZero() {
			payload.CreatedAt = time.Now()
		}

		status := models.TransferInitiationStatusWaitingForValidation
		if payload.Validated {
			status = models.TransferInitiationStatusValidated
		}

		isInstalled, _ := repo.IsInstalled(r.Context(), models.MustConnectorProviderFromString(payload.Provider))
		if !isInstalled {
			handleValidationError(w, r, fmt.Errorf("provider %s is not installed", payload.Provider))

			return
		}

		if payload.SourceAccountID != "" {
			_, err := repo.GetAccount(r.Context(), payload.SourceAccountID)
			if err != nil {
				handleStorageErrors(w, r, fmt.Errorf("failed to get source account: %w", err))

				return
			}
		}

		_, err := repo.GetAccount(r.Context(), payload.DestinationAccountID)
		if err != nil {
			handleStorageErrors(w, r, fmt.Errorf("failed to get destination account: %w", err))

			return
		}

		provider := models.MustConnectorProviderFromString(payload.Provider)
		tf := &models.TransferInitiation{
			ID: models.TransferInitiationID{
				Reference: payload.Reference,
				Provider:  provider,
			},
			CreatedAt:            payload.CreatedAt,
			UpdatedAt:            payload.CreatedAt, // When created, should be the same
			Description:          payload.Description,
			DestinationAccountID: models.MustAccountIDFromString(payload.DestinationAccountID),
			Provider:             provider,
			Type:                 models.MustTransferInitiationTypeFromString(payload.Type),
			Amount:               payload.Amount,
			Asset:                models.Asset(payload.Asset),
			Status:               status,
		}

		if payload.SourceAccountID != "" {
			tf.SourceAccountID = models.MustAccountIDFromString(payload.SourceAccountID)
		}

		if err := repo.CreateTransferInitiation(r.Context(), tf); err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		if status == models.TransferInitiationStatusValidated {
			f, ok := paymentHandlers[provider]
			if !ok {
				handleServerError(w, r, errors.New("no payment handler for provider "+provider.String()))

				return
			}

			err = f(r.Context(), tf)
			if err != nil {
				handleServerError(w, r, err)

				return
			}
		}

		data := &transferInitiationResponse{
			ID:                   tf.ID.String(),
			CreatedAt:            tf.CreatedAt,
			UpdatedAt:            tf.UpdatedAt,
			Description:          tf.Description,
			SourceAccountID:      tf.SourceAccountID.String(),
			DestinationAccountID: tf.DestinationAccountID.String(),
			Provider:             tf.Provider.String(),
			Type:                 tf.Type.String(),
			Amount:               tf.Amount,
			Asset:                tf.Asset.String(),
			Status:               tf.Status.String(),
			Error:                tf.Error,
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[transferInitiationResponse]{
			Data: data,
		})
		if err != nil {
			handleServerError(w, r, err)

			return
		}
	}
}

type udateTransferInitiationStatusRepository interface {
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
}

type updateTransferInitiationStatusRequest struct {
	Status string `json:"status"`
}

func updateTransferInitiationStatusHandler(
	repo udateTransferInitiationStatusRepository,
	paymentHandlers map[models.ConnectorProvider]paymentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &updateTransferInitiationStatusRequest{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			handleValidationError(w, r, err)

			return
		}

		status, err := models.TransferInitiationStatusFromString(payload.Status)
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		switch status {
		case models.TransferInitiationStatusWaitingForValidation:
			handleValidationError(w, r, errors.New("cannot set back transfer initiation status to waiting for validation"))

			return
		case models.TransferInitiationStatusFailed,
			models.TransferInitiationStatusProcessed,
			models.TransferInitiationStatusProcessing:
			handleValidationError(w, r, errors.New("Either VALIDATED or REJECTED status can be set"))

			return
		default:
		}

		transferID, err := models.TransferInitiationIDFromString((mux.Vars(r)["transferID"]))
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		previousTransferInitiation, err := repo.ReadTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		if previousTransferInitiation.Status != models.TransferInitiationStatusWaitingForValidation {
			handleValidationError(w, r, errors.New("only waiting for validation transfer initiation can be updated"))

			return
		}
		previousTransferInitiation.Status = status
		previousTransferInitiation.Attempts++

		err = repo.UpdateTransferInitiationPaymentsStatus(r.Context(), transferID, nil, status, "", previousTransferInitiation.Attempts, time.Now())
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		if status == models.TransferInitiationStatusValidated {
			f, ok := paymentHandlers[previousTransferInitiation.Provider]
			if !ok {
				handleServerError(w, r, errors.New("no payment handler for provider "+previousTransferInitiation.Provider.String()))

				return
			}

			err = f(r.Context(), previousTransferInitiation)
			if err != nil {
				handleServerError(w, r, err)

				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

type retryTransferInitiationRepository interface {
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	UpdateTransferInitiationPaymentsStatus(ctx context.Context, id models.TransferInitiationID, paymentID *models.PaymentID, status models.TransferInitiationStatus, errorMessage string, attempts int, updatedAt time.Time) error
}

func retryTransferInitiationHandler(
	repo retryTransferInitiationRepository,
	paymentHandlers map[models.ConnectorProvider]paymentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transferID, err := models.TransferInitiationIDFromString((mux.Vars(r)["transferID"]))
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		previousTransferInitiation, err := repo.ReadTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		if previousTransferInitiation.Status != models.TransferInitiationStatusFailed {
			handleValidationError(w, r, errors.New("only failed transfer initiation can be updated"))

			return
		}
		previousTransferInitiation.Status = models.TransferInitiationStatusProcessing
		previousTransferInitiation.Attempts++

		err = repo.UpdateTransferInitiationPaymentsStatus(r.Context(), transferID, nil, models.TransferInitiationStatusProcessing, "", previousTransferInitiation.Attempts, time.Now())
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		f, ok := paymentHandlers[previousTransferInitiation.Provider]
		if !ok {
			handleServerError(w, r, errors.New("no payment handler for provider "+previousTransferInitiation.Provider.String()))

			return
		}

		err = f(r.Context(), previousTransferInitiation)
		if err != nil {
			handleServerError(w, r, err)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
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
	ListTransferInitiations(ctx context.Context, pagination storage.Paginator) ([]*models.TransferInitiation, storage.PaginationDetails, error)
}

func listTransferInitiationsHandler(repo listTransferInitiationsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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

		pagination, err := storage.Paginate(pageSize, r.URL.Query().Get("cursor"), sorter)
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

type deleteTransferInitiationRepository interface {
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	DeleteTransferInitiation(ctx context.Context, id models.TransferInitiationID) error
}

func deleteTransferInitiationHandler(repo deleteTransferInitiationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transferID, err := models.TransferInitiationIDFromString(mux.Vars(r)["transferID"])
		if err != nil {
			handleValidationError(w, r, err)

			return
		}

		tf, err := repo.ReadTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		if tf.Status != models.TransferInitiationStatusWaitingForValidation {
			handleValidationError(w, r, errors.New("cannot delete transfer initiation not waiting for validation"))

		}

		err = repo.DeleteTransferInitiation(r.Context(), transferID)
		if err != nil {
			handleStorageErrors(w, r, err)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

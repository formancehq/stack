package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
)

type accountResponse struct {
	ID              string            `json:"id"`
	Reference       string            `json:"reference"`
	CreatedAt       time.Time         `json:"createdAt"`
	ConnectorID     string            `json:"connectorID"`
	Provider        string            `json:"provider"`
	DefaultCurrency string            `json:"defaultCurrency"` // Deprecated: should be removed soon
	DefaultAsset    string            `json:"defaultAsset"`
	AccountName     string            `json:"accountName"`
	Type            string            `json:"type"`
	Metadata        map[string]string `json:"metadata"`
	Pools           []uuid.UUID       `json:"pools"`
	Raw             interface{}       `json:"raw"`
}

func createAccountHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "createAccountHandler")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")

		var req service.CreateAccountRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		span.SetAttributes(
			attribute.String("request.reference", req.Reference),
			attribute.String("request.type", req.Type),
			attribute.String("request.connectorID", req.ConnectorID),
			attribute.String("request.createdAt", req.CreatedAt.String()),
			attribute.String("request.accountName", req.AccountName),
			attribute.String("request.defaultAsset", req.DefaultAsset),
		)

		if err := req.Validate(); err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		account, err := b.GetService().CreateAccount(ctx, &req)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := &accountResponse{
			ID:              account.ID.String(),
			Reference:       account.Reference,
			CreatedAt:       account.CreatedAt,
			ConnectorID:     account.ConnectorID.String(),
			Provider:        account.ConnectorID.Provider.String(),
			DefaultCurrency: account.DefaultAsset.String(),
			DefaultAsset:    account.DefaultAsset.String(),
			AccountName:     account.AccountName,
			Type:            account.Type.String(),
			Raw:             account.RawData,
		}

		if account.Metadata != nil {
			metadata := make(map[string]string)
			for k, v := range account.Metadata {
				metadata[k] = v
			}
			data.Metadata = metadata
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[accountResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func listAccountsHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "listAccountsHandler")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")

		query, err := bunpaginate.Extract[storage.ListAccountsQuery](r, func() (*storage.ListAccountsQuery, error) {
			options, err := getPagination(r, storage.AccountQuery{})
			if err != nil {
				return nil, err
			}
			return pointer.For(storage.NewListAccountsQuery(*options)), nil
		})
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := b.GetService().ListAccounts(ctx, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		ret := cursor.Data
		data := make([]*accountResponse, len(ret))

		for i := range ret {
			accountType := ret[i].Type
			if accountType == models.AccountTypeExternalFormance {
				accountType = models.AccountTypeExternal
			}

			data[i] = &accountResponse{
				ID:              ret[i].ID.String(),
				Reference:       ret[i].Reference,
				CreatedAt:       ret[i].CreatedAt,
				ConnectorID:     ret[i].ConnectorID.String(),
				Provider:        ret[i].ConnectorID.Provider.String(),
				DefaultCurrency: ret[i].DefaultAsset.String(),
				DefaultAsset:    ret[i].DefaultAsset.String(),
				AccountName:     ret[i].AccountName,
				Type:            accountType.String(),
				Raw:             ret[i].RawData,
			}

			if ret[i].Metadata != nil {
				metadata := make(map[string]string)
				for k, v := range ret[i].Metadata {
					metadata[k] = v
				}
				data[i].Metadata = metadata
			}

			data[i].Pools = make([]uuid.UUID, len(ret[i].PoolAccounts))
			for j := range ret[i].PoolAccounts {
				data[i].Pools[j] = ret[i].PoolAccounts[j].PoolID
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*accountResponse]{
			Cursor: &bunpaginate.Cursor[*accountResponse]{
				PageSize: cursor.PageSize,
				HasMore:  cursor.HasMore,
				Previous: cursor.Previous,
				Next:     cursor.Next,
				Data:     data,
			},
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

func readAccountHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "readAccountHandler")
		defer span.End()

		w.Header().Set("Content-Type", "application/json")

		accountID := mux.Vars(r)["accountID"]

		span.SetAttributes(attribute.String("request.accountID", accountID))

		account, err := b.GetService().GetAccount(ctx, accountID)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		accountType := account.Type
		if accountType == models.AccountTypeExternalFormance {
			accountType = models.AccountTypeExternal
		}

		data := &accountResponse{
			ID:              account.ID.String(),
			Reference:       account.Reference,
			CreatedAt:       account.CreatedAt,
			ConnectorID:     account.ConnectorID.String(),
			Provider:        account.ConnectorID.Provider.String(),
			DefaultCurrency: account.DefaultAsset.String(),
			DefaultAsset:    account.DefaultAsset.String(),
			AccountName:     account.AccountName,
			Type:            accountType.String(),
			Raw:             account.RawData,
		}

		if account.Metadata != nil {
			metadata := make(map[string]string)
			for k, v := range account.Metadata {
				metadata[k] = v
			}
			data.Metadata = metadata
		}

		data.Pools = make([]uuid.UUID, len(account.PoolAccounts))
		for j := range account.PoolAccounts {
			data.Pools[j] = account.PoolAccounts[j].PoolID
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[accountResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}

	}
}

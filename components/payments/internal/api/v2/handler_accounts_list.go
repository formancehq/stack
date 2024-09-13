package v2

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
)

// NOTE: in order to maintain previous version compatibility, we need to keep the
// same response structure as the previous version of the API
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
	// TODO(polo): add pools
	// Pools           []uuid.UUID       `json:"pools"`
	Raw interface{} `json:"raw"`
}

func accountsList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_accountsList")
		defer span.End()

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

		cursor, err := backend.AccountsList(ctx, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]*accountResponse, len(cursor.Data))
		for i := range cursor.Data {
			data[i] = &accountResponse{
				ID:          cursor.Data[i].ID.String(),
				Reference:   cursor.Data[i].Reference,
				CreatedAt:   cursor.Data[i].CreatedAt,
				ConnectorID: cursor.Data[i].ConnectorID.String(),
				Provider:    cursor.Data[i].ConnectorID.Provider,
				Type:        string(cursor.Data[i].Type),
				Metadata:    cursor.Data[i].Metadata,
				Raw:         cursor.Data[i].Raw,
			}

			if cursor.Data[i].DefaultAsset != nil {
				data[i].DefaultCurrency = *cursor.Data[i].DefaultAsset
				data[i].DefaultAsset = *cursor.Data[i].DefaultAsset
			}

			if cursor.Data[i].Name != nil {
				data[i].AccountName = *cursor.Data[i].Name
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

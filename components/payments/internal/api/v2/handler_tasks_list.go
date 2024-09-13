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

type listTasksResponseElement struct {
	ID          string          `json:"id"`
	ConnectorID string          `json:"connectorID"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
	Descriptor  json.RawMessage `json:"descriptor"`
	Status      string          `json:"status"`
	State       json.RawMessage `json:"state"`
	Error       string          `json:"error"`
}

func tasksList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_tasksList")
		defer span.End()

		query, err := bunpaginate.Extract[storage.ListSchedulesQuery](r, func() (*storage.ListSchedulesQuery, error) {
			pageSize, err := bunpaginate.GetPageSize(r)
			if err != nil {
				return nil, err
			}

			return pointer.For(storage.NewListSchedulesQuery(bunpaginate.NewPaginatedQueryOptions(storage.ScheduleQuery{}).WithPageSize(pageSize))), nil
		})
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := backend.SchedulesList(ctx, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]listTasksResponseElement, len(cursor.Data))
		for i := range cursor.Data {
			raw, err := json.Marshal(&cursor.Data[i])
			if err != nil {
				otel.RecordError(span, err)
				api.InternalServerError(w, r, err)
				return
			}

			data[i] = listTasksResponseElement{
				ID:          cursor.Data[i].ID,
				ConnectorID: cursor.Data[i].ConnectorID.String(),
				CreatedAt:   cursor.Data[i].CreatedAt.Format(time.RFC3339),
				UpdatedAt:   cursor.Data[i].CreatedAt.Format(time.RFC3339),
				Descriptor:  raw,
				Status:      "ACTIVE",
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[listTasksResponseElement]{
			Cursor: &bunpaginate.Cursor[listTasksResponseElement]{
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

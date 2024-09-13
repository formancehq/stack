package v2

import (
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
)

type paymentResponse struct {
	ID                   string              `json:"id"`
	Reference            string              `json:"reference"`
	SourceAccountID      string              `json:"sourceAccountID"`
	DestinationAccountID string              `json:"destinationAccountID"`
	Type                 string              `json:"type"`
	Provider             string              `json:"provider"`
	ConnectorID          string              `json:"connectorID"`
	Status               string              `json:"status"`
	Amount               *big.Int            `json:"amount"`
	InitialAmount        *big.Int            `json:"initialAmount"`
	Scheme               string              `json:"scheme"`
	Asset                string              `json:"asset"`
	CreatedAt            time.Time           `json:"createdAt"`
	Raw                  interface{}         `json:"raw"`
	Adjustments          []paymentAdjustment `json:"adjustments"`
	Metadata             map[string]string   `json:"metadata"`
}

type paymentAdjustment struct {
	Reference string      `json:"reference" bson:"reference"`
	CreatedAt time.Time   `json:"createdAt" bson:"createdAt"`
	Status    string      `json:"status" bson:"status"`
	Amount    *big.Int    `json:"amount" bson:"amount"`
	Raw       interface{} `json:"raw" bson:"raw"`
}

func paymentsList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_paymentsList")
		defer span.End()

		query, err := bunpaginate.Extract[storage.ListPaymentsQuery](r, func() (*storage.ListPaymentsQuery, error) {
			options, err := getPagination(r, storage.PaymentQuery{})
			if err != nil {
				otel.RecordError(span, err)
				return nil, err
			}
			return pointer.For(storage.NewListPaymentsQuery(*options)), nil
		})
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := backend.PaymentsList(ctx, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]*paymentResponse, len(cursor.Data))
		for i := range cursor.Data {
			data[i] = &paymentResponse{
				ID:            cursor.Data[i].ID.String(),
				Reference:     cursor.Data[i].Reference,
				Type:          cursor.Data[i].Type.String(),
				Provider:      cursor.Data[i].ConnectorID.Provider,
				ConnectorID:   cursor.Data[i].ConnectorID.String(),
				Status:        cursor.Data[i].Status.String(),
				Amount:        cursor.Data[i].Amount,
				InitialAmount: cursor.Data[i].InitialAmount,
				Scheme:        cursor.Data[i].Scheme.String(),
				Asset:         cursor.Data[i].Asset,
				CreatedAt:     cursor.Data[i].CreatedAt,
				Adjustments:   []paymentAdjustment{},
				Metadata:      cursor.Data[i].Metadata,
			}

			if cursor.Data[i].SourceAccountID != nil {
				data[i].SourceAccountID = cursor.Data[i].SourceAccountID.String()
			}

			if cursor.Data[i].DestinationAccountID != nil {
				data[i].DestinationAccountID = cursor.Data[i].DestinationAccountID.String()
			}

			data[i].Adjustments = make([]paymentAdjustment, len(cursor.Data[i].Adjustments))
			for j := range cursor.Data[i].Adjustments {
				data[i].Adjustments[j] = paymentAdjustment{
					Reference: cursor.Data[i].Adjustments[j].ID.Reference,
					CreatedAt: cursor.Data[i].Adjustments[j].CreatedAt,
					Status:    cursor.Data[i].Adjustments[j].Status.String(),
					Amount:    cursor.Data[i].Adjustments[j].Amount,
					Raw:       cursor.Data[i].Adjustments[j].Raw,
				}
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*paymentResponse]{
			Cursor: &bunpaginate.Cursor[*paymentResponse]{
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

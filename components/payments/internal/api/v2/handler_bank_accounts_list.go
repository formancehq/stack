package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/payments/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
)

func bankAccountsList(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_bankAccountsList")
		defer span.End()

		query, err := bunpaginate.Extract[storage.ListBankAccountsQuery](r, func() (*storage.ListBankAccountsQuery, error) {
			options, err := getPagination(r, storage.BankAccountQuery{})
			if err != nil {
				otel.RecordError(span, err)
				return nil, err
			}
			return pointer.For(storage.NewListBankAccountsQuery(*options)), nil
		})
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrValidation, err)
			return
		}

		cursor, err := backend.BankAccountsList(ctx, *query)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]*bankAccountResponse, len(cursor.Data))
		for i := range cursor.Data {
			data[i] = &bankAccountResponse{
				ID:        cursor.Data[i].ID.String(),
				Name:      cursor.Data[i].Name,
				CreatedAt: cursor.Data[i].CreatedAt,
				Metadata:  cursor.Data[i].Metadata,
			}

			if cursor.Data[i].IBAN != nil {
				data[i].Iban = *cursor.Data[i].IBAN
			}

			if cursor.Data[i].AccountNumber != nil {
				data[i].AccountNumber = *cursor.Data[i].AccountNumber
			}

			if cursor.Data[i].SwiftBicCode != nil {
				data[i].SwiftBicCode = *cursor.Data[i].SwiftBicCode
			}

			if cursor.Data[i].Country != nil {
				data[i].Country = *cursor.Data[i].Country
			}

			data[i].RelatedAccounts = make([]*bankAccountRelatedAccountsResponse, len(cursor.Data[i].RelatedAccounts))
			for j := range cursor.Data[i].RelatedAccounts {
				data[i].RelatedAccounts[j] = &bankAccountRelatedAccountsResponse{
					ID:          "",
					CreatedAt:   cursor.Data[i].RelatedAccounts[j].CreatedAt,
					AccountID:   cursor.Data[i].RelatedAccounts[j].AccountID.String(),
					ConnectorID: cursor.Data[i].RelatedAccounts[j].ConnectorID.String(),
					Provider:    cursor.Data[i].RelatedAccounts[j].ConnectorID.Provider,
				}
			}
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[*bankAccountResponse]{
			Cursor: &bunpaginate.Cursor[*bankAccountResponse]{
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

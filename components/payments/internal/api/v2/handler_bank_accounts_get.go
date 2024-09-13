package v2

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

func bankAccountsGet(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_bankAccountsGet")
		defer span.End()

		id, err := uuid.Parse(bankAccountID(r))
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		bankAccount, err := backend.BankAccountsGet(ctx, id)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		data := &bankAccountResponse{
			ID:        bankAccount.ID.String(),
			Name:      bankAccount.Name,
			CreatedAt: bankAccount.CreatedAt,
			Metadata:  bankAccount.Metadata,
		}

		if bankAccount.IBAN != nil {
			data.Iban = *bankAccount.IBAN
		}

		if bankAccount.AccountNumber != nil {
			data.AccountNumber = *bankAccount.AccountNumber
		}

		if bankAccount.SwiftBicCode != nil {
			data.SwiftBicCode = *bankAccount.SwiftBicCode
		}

		if bankAccount.Country != nil {
			data.Country = *bankAccount.Country
		}

		for _, relatedAccount := range bankAccount.RelatedAccounts {
			data.RelatedAccounts = append(data.RelatedAccounts, &bankAccountRelatedAccountsResponse{
				ID:          "",
				CreatedAt:   relatedAccount.CreatedAt,
				AccountID:   relatedAccount.AccountID.String(),
				ConnectorID: relatedAccount.ConnectorID.String(),
				Provider:    relatedAccount.ConnectorID.Provider,
			})
		}

		err = json.NewEncoder(w).Encode(api.BaseResponse[bankAccountResponse]{
			Data: data,
		})
		if err != nil {
			otel.RecordError(span, err)
			api.InternalServerError(w, r, err)
			return
		}
	}
}

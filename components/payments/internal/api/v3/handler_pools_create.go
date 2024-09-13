package v3

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/google/uuid"
)

type createPoolRequest struct {
	Name       string   `json:"name"`
	AccountIDs []string `json:"accountIDs"`
}

func poolsCreate(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v3_poolsCreate")
		defer span.End()

		var createPoolRequest createPoolRequest
		err := json.NewDecoder(r.Body).Decode(&createPoolRequest)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		pool := models.Pool{
			ID:        uuid.New(),
			Name:      createPoolRequest.Name,
			CreatedAt: time.Now().UTC(),
		}

		accounts := make([]models.PoolAccounts, len(createPoolRequest.AccountIDs))
		for i, accountID := range createPoolRequest.AccountIDs {
			aID, err := models.AccountIDFromString(accountID)
			if err != nil {
				otel.RecordError(span, err)
				api.BadRequest(w, ErrValidation, err)
				return
			}

			accounts[i] = models.PoolAccounts{
				PoolID:    pool.ID,
				AccountID: aID,
			}
		}
		pool.PoolAccounts = accounts

		err = backend.PoolsCreate(ctx, pool)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.Created(w, pool.ID.String())
	}
}

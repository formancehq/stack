package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pkg/errors"

	"github.com/gorilla/mux"
)

type updateMetadataRepository interface {
	UpdatePaymentMetadata(ctx context.Context, paymentID models.PaymentID, metadata map[string]string) error
}

type updateMetadataRequest map[string]string

func updateMetadataHandler(repo updateMetadataRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID, err := models.PaymentIDFromString(mux.Vars(r)["paymentID"])
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		var metadata updateMetadataRequest

		if r.ContentLength == 0 {
			api.BadRequest(w, ErrMissingBody, errors.New("body is required"))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&metadata)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		err = repo.UpdatePaymentMetadata(r.Context(), *paymentID, metadata)
		if err != nil {
			api.InternalServerError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

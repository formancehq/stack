package api

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/pkg/errors"

	"github.com/gorilla/mux"
)

func updateMetadataHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID, err := models.PaymentIDFromString(mux.Vars(r)["paymentID"])
		if err != nil {
			api.BadRequest(w, ErrInvalidID, err)
			return
		}

		var metadata service.UpdateMetadataRequest
		if r.ContentLength == 0 {
			api.BadRequest(w, ErrMissingOrInvalidBody, errors.New("body is required"))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&metadata)
		if err != nil {
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		err = b.GetService().UpdatePaymentMetadata(r.Context(), *paymentID, metadata)
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

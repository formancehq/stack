package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

type updateMetadataRepository interface {
	UpdatePaymentMetadata(ctx context.Context, paymentID uuid.UUID, metadata map[string]string) error
}

type updateMetadataRequest map[string]string

func updateMetadataHandler(repo updateMetadataRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID, err := uuid.Parse(mux.Vars(r)["paymentID"])
		if err != nil {
			handleErrorBadRequest(w, r, err)

			return
		}

		var metadata updateMetadataRequest

		if r.ContentLength == 0 {
			handleErrorBadRequest(w, r, errors.New("body is required"))

			return
		}

		err = json.NewDecoder(r.Body).Decode(&metadata)
		if err != nil {
			handleError(w, r, err)

			return
		}

		err = repo.UpdatePaymentMetadata(r.Context(), paymentID, metadata)
		if err != nil {
			handleError(w, r, err)

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

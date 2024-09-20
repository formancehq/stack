package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"

	"github.com/uptrace/bun"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/logging"
	"go.opentelemetry.io/otel/trace"
)

func validationError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(api.ErrorResponse{
		ErrorCode:    "VALIDATION",
		ErrorMessage: err.Error(),
	}); err != nil {
		logging.FromContext(r.Context()).Info("Error validating request: %s", err)
	}
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(api.ErrorResponse{
		ErrorCode:    "INTERNAL",
		ErrorMessage: err.Error(),
	}); err != nil {
		trace.SpanFromContext(r.Context()).RecordError(err)
	}
}

func writeJSONObject[T any](w http.ResponseWriter, r *http.Request, v T) {
	if err := json.NewEncoder(w).Encode(api.BaseResponse[T]{
		Data: &v,
	}); err != nil {
		trace.SpanFromContext(r.Context()).RecordError(err)
	}
}

func writeCreatedJSONObject(w http.ResponseWriter, r *http.Request, v any, id string) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "./"+id)

	writeJSONObject(w, r, v)
}

func readJSONObject[T any](w http.ResponseWriter, r *http.Request) *T {
	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		validationError(w, r, err)
		return nil
	}
	return &t
}

func findById[T any](w http.ResponseWriter, r *http.Request, db *bun.DB, params string) T {
	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	err := db.NewSelect().
		Model(t).
		Limit(1).
		Where("id = ?", chi.URLParam(r, params)).
		Scan(r.Context())
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
		default:
			internalServerError(w, r, err)
		}
		var zeroValue T
		return zeroValue
	}
	return t
}

func createObject(w http.ResponseWriter, r *http.Request, db *bun.DB, v any) error {
	_, err := db.NewInsert().Model(v).Exec(r.Context())
	if err != nil {
		internalServerError(w, r, err)
	}
	return err
}

func mapList[I any, O any](items []I, fn func(I) O) []O {
	ret := make([]O, 0)
	for _, item := range items {
		ret = append(ret, fn(item))
	}
	return ret
}

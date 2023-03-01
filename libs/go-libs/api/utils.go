package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

const defaultLimit = 15

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter, code string, err error) {
	w.WriteHeader(http.StatusBadRequest)
	writeJSON(w, ErrorResponse{
		ErrorCode:    code,
		ErrorMessage: err.Error(),
	})
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	logging.FromContext(r.Context()).Error(err)

	w.WriteHeader(http.StatusInternalServerError)
	writeJSON(w, ErrorResponse{
		ErrorCode:    "INTERNAL_ERROR",
		ErrorMessage: err.Error(),
	})
}

func Created(w http.ResponseWriter, v any) {
	w.WriteHeader(http.StatusCreated)
	Ok(w, v)
}

func RawOk(w http.ResponseWriter, v any) {
	writeJSON(w, v)
}

func Ok(w http.ResponseWriter, v any) {
	RawOk(w, BaseResponse[any]{
		Data: &v,
	})
}

func RenderCursor[T any](w http.ResponseWriter, v Cursor[T]) {
	writeJSON(w, BaseResponse[T]{
		Cursor: &v,
	})
}

func CursorFromListResponse[T any, V any](w http.ResponseWriter, query ListQuery[V], response *ListResponse[T]) {
	RenderCursor(w, Cursor[T]{
		PageSize: query.Limit,
		HasMore:  response.HasMore,
		Previous: response.Previous,
		Next:     response.Next,
		Data:     response.Data,
	})
}

func ParsePaginationToken(r *http.Request) string {
	return r.URL.Query().Get("RenderCursor")
}

func ParsePageSize(r *http.Request) int {
	pageSize := r.URL.Query().Get("pageSize")
	if pageSize == "" {
		return defaultLimit
	}

	v, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(v)
}

func ReadPaginatedRequest[T any](r *http.Request, f func(r *http.Request) T) ListQuery[T] {
	var payload T
	if f != nil {
		payload = f(r)
	}
	return ListQuery[T]{
		Pagination: Pagination{
			Limit:           ParsePageSize(r),
			PaginationToken: ParsePaginationToken(r),
		},
		Payload: payload,
	}
}

func GetQueryMap(m map[string][]string, key string) map[string]string {
	dicts := make(map[string]string)
	for k, v := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				dicts[k[i+1:][:j]] = v[0]
			}
		}
	}
	return dicts
}

type ListResponse[T any] struct {
	Data           []T
	Next, Previous string
	HasMore        bool
}

type Pagination struct {
	Limit           int
	PaginationToken string
}

type ListQuery[T any] struct {
	Pagination
	Payload T
}

type Mapper[SRC any, DST any] func(src SRC) DST

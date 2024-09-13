package v2

import (
	"io"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
)

func getQueryBuilder(r *http.Request) (query.Builder, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		return query.ParseJSON(string(data))
	} else {
		// In order to be backward compatible
		return query.ParseJSON(r.URL.Query().Get("query"))
	}
}

func getPagination[T any](r *http.Request, options T) (*bunpaginate.PaginatedQueryOptions[T], error) {
	qb, err := getQueryBuilder(r)
	if err != nil {
		return nil, err
	}

	pageSize, err := bunpaginate.GetPageSize(r)
	if err != nil {
		return nil, err
	}

	return pointer.For(bunpaginate.NewPaginatedQueryOptions(options).WithQueryBuilder(qb).WithPageSize(pageSize)), nil
}

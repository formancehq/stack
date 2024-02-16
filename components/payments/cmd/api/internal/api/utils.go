package api

import (
	"net/http"
	"strings"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/pkg/errors"
)

func getQueryBuilder(r *http.Request) (query.Builder, error) {
	return query.ParseJSON(r.URL.Query().Get("query"))
}

func getSorter(r *http.Request) (storage.Sorter, error) {
	var sorter storage.Sorter

	if sortParams := r.URL.Query()["sort"]; sortParams != nil {
		for _, s := range sortParams {
			parts := strings.SplitN(s, ":", 2)

			var order storage.SortOrder

			if len(parts) > 1 {
				//nolint:goconst // allow duplicate string
				switch parts[1] {
				case "asc", "ASC":
					order = storage.SortOrderAsc
				case "dsc", "desc", "DSC", "DESC":
					order = storage.SortOrderDesc
				default:
					return sorter, errors.New("sort order not well specified, got " + parts[1])
				}
			}

			column := parts[0]

			sorter = sorter.Add(column, order)
		}
	}

	return sorter, nil
}

func getPagination[T any](r *http.Request, options T) (*storage.PaginatedQueryOptions[T], error) {
	qb, err := getQueryBuilder(r)
	if err != nil {
		return nil, err
	}

	sorter, err := getSorter(r)
	if err != nil {
		return nil, err
	}

	pageSize, err := bunpaginate.GetPageSize(r)
	if err != nil {
		return nil, err
	}

	return pointer.For(storage.NewPaginatedQueryOptions(options).WithQueryBuilder(qb).WithSorter(sorter).WithPageSize(pageSize)), nil
}

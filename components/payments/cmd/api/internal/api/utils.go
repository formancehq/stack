package api

import (
	"net/http"
	"strings"

	"github.com/formancehq/payments/cmd/api/internal/storage"
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

func getPagination(r *http.Request) (storage.PaginatorQuery, error) {
	qb, err := getQueryBuilder(r)
	if err != nil {
		return storage.PaginatorQuery{}, err
	}

	sorter, err := getSorter(r)
	if err != nil {
		return storage.PaginatorQuery{}, err
	}

	pageSize, err := pageSizeQueryParam(r)
	if err != nil {
		return storage.PaginatorQuery{}, err
	}

	pagination, err := storage.Paginate(pageSize, r.URL.Query().Get("cursor"), sorter, qb)
	if err != nil {
		return storage.PaginatorQuery{}, err
	}

	return pagination, nil
}

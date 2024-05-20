package v2

import (
	"io"
	"net/http"

	storage "github.com/formancehq/reconciliation/internal/storage/v2"
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
	}
	return nil, nil
}

func getPaginatedQueryOptionsRules(r *http.Request) (*bunpaginate.PaginatedQueryOptions[storage.RulesFilters], error) {
	qb, err := getQueryBuilder(r)
	if err != nil {
		return nil, err
	}

	pageSize, err := getPageSize(r)
	if err != nil {
		return nil, err
	}

	filters := storage.RulesFilters{}
	return pointer.For(bunpaginate.NewPaginatedQueryOptions(filters).
		WithQueryBuilder(qb).
		WithPageSize(pageSize)), nil
}

func getPaginatedQueryOptionsPolicies(r *http.Request) (*bunpaginate.PaginatedQueryOptions[storage.PoliciesFilters], error) {
	qb, err := getQueryBuilder(r)
	if err != nil {
		return nil, err
	}

	pageSize, err := getPageSize(r)
	if err != nil {
		return nil, err
	}

	filters := storage.PoliciesFilters{}
	return pointer.For(bunpaginate.NewPaginatedQueryOptions(filters).
		WithQueryBuilder(qb).
		WithPageSize(pageSize)), nil
}

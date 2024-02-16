package storage

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

type PaginatedQueryOptions[T any] struct {
	QueryBuilder query.Builder `json:"qb"`
	Sorter       Sorter
	PageSize     uint64 `json:"pageSize"`
	Options      T      `json:"options"`
}

func (opts PaginatedQueryOptions[T]) WithQueryBuilder(qb query.Builder) PaginatedQueryOptions[T] {
	opts.QueryBuilder = qb

	return opts
}

func (opts PaginatedQueryOptions[T]) WithSorter(sorter Sorter) PaginatedQueryOptions[T] {
	opts.Sorter = sorter

	return opts
}

func (opts PaginatedQueryOptions[T]) WithPageSize(pageSize uint64) PaginatedQueryOptions[T] {
	opts.PageSize = pageSize

	return opts
}

func NewPaginatedQueryOptions[T any](options T) PaginatedQueryOptions[T] {
	return PaginatedQueryOptions[T]{
		Options:  options,
		PageSize: bunpaginate.QueryDefaultPageSize,
	}
}

func PaginateWithOffset[FILTERS any, RETURN any](s *Storage, ctx context.Context,
	q *bunpaginate.OffsetPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*api.Cursor[RETURN], error) {
	query := s.db.NewSelect()
	return bunpaginate.UsingOffset[FILTERS, RETURN](ctx, query, *q, builders...)
}

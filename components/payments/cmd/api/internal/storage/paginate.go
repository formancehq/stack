package storage

import (
	"context"
	"encoding/json"

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

func (v *PaginatedQueryOptions[T]) UnmarshalJSON(data []byte) error {
	type aux struct {
		QueryBuilder json.RawMessage `json:"qb"`
		Sorter       Sorter          `json:"Sorter"`
		PageSize     uint64          `json:"pageSize"`
		Options      T               `json:"options"`
	}
	x := &aux{}
	if err := json.Unmarshal(data, x); err != nil {
		return err
	}

	*v = PaginatedQueryOptions[T]{
		PageSize: x.PageSize,
		Options:  x.Options,
		Sorter:   x.Sorter,
	}

	var err error
	if x.QueryBuilder != nil {
		v.QueryBuilder, err = query.ParseJSON(string(x.QueryBuilder))
		if err != nil {
			return err
		}
	}

	return nil
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
	q *bunpaginate.OffsetPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*bunpaginate.Cursor[RETURN], error) {
	query := s.db.NewSelect()
	return bunpaginate.UsingOffset[FILTERS, RETURN](ctx, query, *q, builders...)
}

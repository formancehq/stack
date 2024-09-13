package storage

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/uptrace/bun"
)

func paginateWithOffset[FILTERS any, RETURN any](s *store, ctx context.Context,
	q *bunpaginate.OffsetPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*bunpaginate.Cursor[RETURN], error) {
	query := s.db.NewSelect()
	return bunpaginate.UsingOffset[FILTERS, RETURN](ctx, query, *q, builders...)
}

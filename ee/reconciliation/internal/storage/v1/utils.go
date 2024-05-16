package v1

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/uptrace/bun"
)

func paginateWithOffset[FILTERS any, RETURN any](s *Storage, ctx context.Context,
	q *bunpaginate.OffsetPaginatedQuery[FILTERS], builders ...func(query *bun.SelectQuery) *bun.SelectQuery) (*bunpaginate.Cursor[RETURN], error) {

	query := s.db.NewSelect()
	for _, builder := range builders {
		query = query.Apply(builder)
	}

	return bunpaginate.UsingOffset[FILTERS, RETURN](ctx, query, *q)
}

package storage

import (
	"github.com/uptrace/bun"
)

type Paginator struct {
	offset int
	limit  int
}

func Paginate(offset, limit int) Paginator {
	return Paginator{offset, limit}
}

func (p Paginator) apply(query *bun.SelectQuery) *bun.SelectQuery {
	return query.Offset(p.offset).Limit(p.limit)
}

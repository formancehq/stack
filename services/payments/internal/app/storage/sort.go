package storage

import (
	"fmt"

	"github.com/uptrace/bun"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type sortExpression struct {
	Column string
	Order  SortOrder
}

type Sorter []sortExpression

func (s Sorter) Add(column string, order SortOrder) Sorter {
	return append(s, sortExpression{column, order})
}

func (s Sorter) apply(query *bun.SelectQuery) *bun.SelectQuery {
	for _, expr := range s {
		query = query.Order(fmt.Sprintf("%s %s", expr.Column, expr.Order))
	}

	return query
}

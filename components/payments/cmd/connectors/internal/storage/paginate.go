package storage

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

const (
	defaultPageSize = 15
	maxPageSize     = 100
)

type PaginationDetails struct {
	PageSize     int
	HasMore      bool
	PreviousPage string
	NextPage     string
}

type baseCursor struct {
	Reference string `json:"reference"`
	Sorter    Sorter `json:"sorter"`
	Next      bool   `json:"next"`
}

func (c baseCursor) Encode() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("error marshaling baseCursor: %w", err)
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

type PaginatorQuery struct {
	pageSize int
	token    string

	queryBuilder query.Builder
	cursor       baseCursor
	sorter       Sorter
}

func Paginate(pageSize int, token string, sorter Sorter, queryBuilder query.Builder) (PaginatorQuery, error) {
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	var cursor baseCursor

	if token != "" {
		tokenBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return PaginatorQuery{}, fmt.Errorf("error decoding token: %w", err)
		}

		err = json.Unmarshal(tokenBytes, &cursor)
		if err != nil {
			return PaginatorQuery{}, fmt.Errorf("error unmarshaling baseCursor: %w", err)
		}
	}

	return PaginatorQuery{pageSize, token, queryBuilder, cursor, sorter}, nil
}

func (p PaginatorQuery) apply(query *bun.SelectQuery, column string) *bun.SelectQuery {
	query = query.Limit(p.pageSize + 1)

	if p.cursor.Reference == "" {
		if p.sorter != nil {
			query = p.sorter.apply(query)
		}

		return query.Order(column + " DESC")
	}

	if p.cursor.Sorter != nil {
		query = p.cursor.Sorter.apply(query)
	}

	if p.cursor.Next {
		return query.Where(fmt.Sprintf("%s < ?", column), p.cursor.Reference).Order(column + " DESC")
	}

	return query.Where(fmt.Sprintf("%s >= ?", column), p.cursor.Reference).Order(column + " ASC")
}

func (p PaginatorQuery) hasPrevious(ctx context.Context, query *bun.SelectQuery, column, reference string) (bool, error) {
	query = query.Limit(1).Order(column + " DESC")

	if p.cursor.Reference == "" {
		if p.sorter != nil {
			query = p.sorter.apply(query)
		}
	}

	if p.cursor.Sorter != nil {
		query = p.cursor.Sorter.apply(query)
	}

	query = query.Where(fmt.Sprintf("%s > ?", column), reference)

	exists, err := query.Exists(ctx)
	if err != nil {
		return false, e("error checking if previous page exists", err)
	}

	return exists, nil
}

func (p PaginatorQuery) paginationDetails(hasMore, hasPrevious bool, firstReference, lastReference string) (PaginationDetails, error) {
	var (
		previousPage string
		nextPage     string
		err          error
	)

	if hasPrevious && firstReference != "" {
		previousPage, err = baseCursor{
			Reference: firstReference,
			Sorter:    p.sorter,
			Next:      false,
		}.Encode()
		if err != nil {
			return PaginationDetails{}, fmt.Errorf("error encoding previous page cursor: %w", err)
		}
	}

	if hasMore && lastReference != "" {
		nextPage, err = baseCursor{
			Reference: lastReference,
			Sorter:    p.sorter,
			Next:      true,
		}.Encode()
		if err != nil {
			return PaginationDetails{}, fmt.Errorf("error encoding next page cursor: %w", err)
		}
	}

	return PaginationDetails{
		PageSize:     p.pageSize,
		HasMore:      hasMore,
		PreviousPage: previousPage,
		NextPage:     nextPage,
	}, nil
}

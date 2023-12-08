package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicateKeyValue = errors.New("duplicate key value")
var ErrInvalidQuery = errors.New("invalid query")

func e(msg string, err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrDuplicateKeyValue
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	return fmt.Errorf("%s: %w", msg, err)
}

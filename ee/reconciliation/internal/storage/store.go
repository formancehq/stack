package storage

import (
	v1 "github.com/formancehq/reconciliation/internal/storage/v1"
	"github.com/uptrace/bun"
)

type Storage struct {
	db *bun.DB

	*v1.Storage
}

func NewStorage(db *bun.DB) *Storage {
	return &Storage{
		db:      db,
		Storage: v1.NewStorage(db),
	}
}

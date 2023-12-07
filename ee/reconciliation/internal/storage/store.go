package storage

import (
	"github.com/uptrace/bun"
)

type Storage struct {
	db *bun.DB
}

func NewStorage(db *bun.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) DB() *bun.DB {
	return s.db
}

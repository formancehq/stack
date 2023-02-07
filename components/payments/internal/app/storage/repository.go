package storage

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
)

type Storage struct {
	db                  *bun.DB
	configEncryptionKey string
}

const encryptionOptions = "compress-algo=1, cipher-algo=aes256"

func newStorage(db *bun.DB, configEncryptionKey string) *Storage {
	return &Storage{db: db, configEncryptionKey: configEncryptionKey}
}

//nolint:unused // used in debug mode
func (s *Storage) debug() {
	s.db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
}

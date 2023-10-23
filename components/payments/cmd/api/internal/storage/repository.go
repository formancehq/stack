package storage

import (
	"github.com/uptrace/bun"
)

type Storage struct {
	db                  *bun.DB
	configEncryptionKey string
}

const encryptionOptions = "compress-algo=1, cipher-algo=aes256"

func newStorage(db *bun.DB, configEncryptionKey string) *Storage {
	return &Storage{db: db, configEncryptionKey: configEncryptionKey}
}

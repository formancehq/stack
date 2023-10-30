package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
)

type Storage struct {
	db                  *bun.DB
	configEncryptionKey string
}

const encryptionOptions = "compress-algo=1, cipher-algo=aes256"

func NewStorage(db *bun.DB, configEncryptionKey string) *Storage {
	return &Storage{db: db, configEncryptionKey: configEncryptionKey}
}

//nolint:unused // used in debug mode
func (s *Storage) debug() {
	s.db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
}

type Reader interface {
	ReadTransferInitiation(ctx context.Context, id models.TransferInitiationID) (*models.TransferInitiation, error)
	GetAccount(ctx context.Context, id string) (*models.Account, error)
}

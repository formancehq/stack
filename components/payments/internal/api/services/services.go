package services

import (
	"github.com/formancehq/payments/internal/connectors/engine"
	"github.com/formancehq/payments/internal/storage"
)

type Service struct {
	storage storage.Storage

	engine engine.Engine
}

func New(storage storage.Storage, engine engine.Engine) *Service {
	return &Service{
		storage: storage,
		engine:  engine,
	}
}

package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Webhook struct {
	bun.BaseModel `bun:"connectors.webhook"`

	ID          uuid.UUID
	ConnectorID ConnectorID
	RequestBody []byte
}

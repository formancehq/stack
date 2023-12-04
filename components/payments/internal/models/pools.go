package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PoolAccounts struct {
	bun.BaseModel `bun:"accounts.pool_accounts"`

	PoolID    uuid.UUID `bun:",pk,notnull"`
	AccountID AccountID `bun:",pk,notnull"`
}

type Pool struct {
	bun.BaseModel `bun:"accounts.pools"`

	ID        uuid.UUID `bun:",pk,notnull"`
	Name      string
	CreatedAt time.Time

	PoolAccounts []*PoolAccounts `bun:"rel:has-many,join:id=pool_id"`
}

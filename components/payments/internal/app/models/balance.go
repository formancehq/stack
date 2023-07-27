package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Balance struct {
	bun.BaseModel `bun:"accounts.balances"`

	AccountID     AccountID `bun:"type:character varying,nullzero"`
	Currency      string
	Balance       int64
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

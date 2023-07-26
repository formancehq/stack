package models

import (
	"math/big"
	"time"

	"github.com/uptrace/bun"
)

type Balance struct {
	bun.BaseModel `bun:"accounts.balances"`

	AccountID     AccountID `bun:"type:character varying,nullzero"`
	Currency      string
	Balance       *big.Int `bun:"type:numeric"`
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

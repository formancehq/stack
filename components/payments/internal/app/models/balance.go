package models

import (
	"math/big"
	"time"

	"github.com/uptrace/bun"
)

type Balance struct {
	bun.BaseModel `bun:"accounts.balances"`

	AccountID     AccountID `bun:"type:character varying,nullzero"`
	Asset         Asset     `bun:"currency"`
	Balance       *big.Int  `bun:"type:numeric"`
	CreatedAt     time.Time
	LastUpdatedAt time.Time
}

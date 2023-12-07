package models

import (
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type ReconciliationStatus string

var (
	ReconciliationNotOK ReconciliationStatus = "NOT_OK"
	ReconciliationOK    ReconciliationStatus = "OK"
)

func (r ReconciliationStatus) String() string {
	return string(r)
}

type Reconciliation struct {
	bun.BaseModel `bun:"reconciliations.reconciliation"`

	ID               uuid.UUID `bun:",pk,nullzero"`
	PolicyID         uuid.UUID `bun:",nullzero"`
	CreatedAt        time.Time `bun:",nullzero"`
	ReconciledAt     time.Time `bun:",nullzero"`
	Status           ReconciliationStatus
	LedgerBalances   map[string]*big.Int `bun:",jsonb"`
	PaymentsBalances map[string]*big.Int `bun:",jsonb"`
	Error            string
}

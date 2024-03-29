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
	bun.BaseModel `bun:"reconciliations.reconciliation" json:"-"`

	ID                   uuid.UUID            `bun:",pk,nullzero" json:"id"`
	PolicyID             uuid.UUID            `bun:",nullzero" json:"policyID"`
	CreatedAt            time.Time            `bun:",nullzero" json:"createdAt"`
	ReconciledAtLedger   time.Time            `bun:",nullzero" json:"reconciledAtLedger"`
	ReconciledAtPayments time.Time            `bun:",nullzero" json:"reconciledAtPayments"`
	Status               ReconciliationStatus `json:"status"`
	LedgerBalances       map[string]*big.Int  `bun:",jsonb" json:"ledgerBalances"`
	PaymentsBalances     map[string]*big.Int  `bun:",jsonb" json:"paymentsBalances"`
	DriftBalances        map[string]*big.Int  `bun:",jsonb" json:"driftBalances"`
	Error                string               `json:"error"`
}

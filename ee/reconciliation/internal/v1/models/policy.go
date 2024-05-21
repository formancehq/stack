package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Policy struct {
	bun.BaseModel `bun:"reconciliations.policy" json:"-"`

	// Policy Related fields
	ID        uuid.UUID `bun:",pk,nullzero" json:"id"`
	CreatedAt time.Time `bun:",notnull" json:"createdAt"`
	Name      string    `bun:",notnull" json:"name"`

	// Reconciliation Needed fields
	LedgerName     string                 `bun:",notnull" json:"ledgerName"`
	LedgerQuery    map[string]interface{} `bun:",type:jsonb,notnull" json:"ledgerQuery"`
	PaymentsPoolID uuid.UUID              `bun:",notnull" json:"paymentsPoolID"`
}

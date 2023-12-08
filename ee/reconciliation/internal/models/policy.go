package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Policy struct {
	bun.BaseModel `bun:"reconciliations.policy"`

	// Policy Related fields
	ID        uuid.UUID `bun:",pk,nullzero"`
	CreatedAt time.Time `bun:",notnull"`
	Name      string    `bun:",notnull"`

	// Reconciliation Needed fields
	LedgerName     string                 `bun:",notnull"`
	LedgerQuery    map[string]interface{} `bun:",type:jsonb,notnull"`
	PaymentsPoolID uuid.UUID              `bun:",notnull"`
}

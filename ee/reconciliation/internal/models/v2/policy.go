package v2

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Policy struct {
	bun.BaseModel `bun:"reconciliationsv2.policy" json:"-"`

	ID        uuid.UUID `bun:",pk,notnull" json:"id"`
	Name      string    `bun:",notnull" json:"name"`
	CreatedAt time.Time `bun:",notnull" json:"createdAt"`
	UpdatedAt time.Time `bun:",notnull" json:"updatedAt"`

	Enabled bool `bun:",notnull" json:"enabled"`

	Rules []string `json:"rules"`
}

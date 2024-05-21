package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PolicyType string

const (
	PolicyTypeAccountBased     = "ACCOUNT_BASED"
	PolicyTypeTransactionBased = "TRANSACTION_BASED"
)

func PolicyTypeFromString(s string) (PolicyType, error) {
	switch s {
	case PolicyTypeAccountBased:
		return PolicyTypeAccountBased, nil
	case PolicyTypeTransactionBased:
		return PolicyTypeTransactionBased, nil
	default:
		return "", fmt.Errorf("unknown policy type %q", s)
	}
}

type Policy struct {
	bun.BaseModel `bun:"reconciliationsv2.policy" json:"-"`

	ID        uuid.UUID  `bun:",pk,notnull" json:"id"`
	Name      string     `bun:",notnull" json:"name"`
	CreatedAt time.Time  `bun:",notnull" json:"createdAt"`
	UpdatedAt time.Time  `bun:",notnull" json:"updatedAt"`
	Type      PolicyType `bun:",notnull" json:"type"`

	Enabled bool `bun:",notnull" json:"enabled"`

	Rules []string `json:"rules"`
}

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
	bun.BaseModel `bun:"reconciliationsv2.policies" json:"-"`

	ID        uuid.UUID  `bun:",pk,notnull" json:"id"`
	Name      string     `bun:",notnull" json:"name"`
	CreatedAt time.Time  `bun:",notnull" json:"createdAt"`
	UpdatedAt time.Time  `bun:",notnull" json:"updatedAt"`
	Type      PolicyType `bun:",notnull" json:"type"`

	// TODO(polo): put it somewhere else, as it has nothing to do with the
	// policy itself.
	// Fine for a first poc
	AdditionalConfig map[string]interface{} `bun:",type:jsonb,notnull" json:"additionalConfig"`

	Enabled bool `bun:",notnull" json:"enabled"`

	Rule uuid.UUID `json:"rule"`
}

func AdditionalConfig(ledgerOverdraftAccountReference string) map[string]interface{} {
	if ledgerOverdraftAccountReference == "" {
		ledgerOverdraftAccountReference = "world"
	}

	return map[string]interface{}{
		"ledgerOverdraftAccountReference": ledgerOverdraftAccountReference,
	}
}

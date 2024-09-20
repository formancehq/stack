package models

import (
	"encoding/base64"
	"time"

	"github.com/gibson042/canonicaljson-go"
	"github.com/google/uuid"
)

type Pool struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`

	PoolAccounts []PoolAccounts `json:"poolAccounts"`
}

func (p *Pool) IdempotencyKey() string {
	relatedAccounts := make([]string, len(p.PoolAccounts))
	for i := range p.PoolAccounts {
		relatedAccounts[i] = p.PoolAccounts[i].AccountID.String()
	}
	var ik = struct {
		ID              string
		RelatedAccounts []string
	}{
		ID:              p.ID.String(),
		RelatedAccounts: relatedAccounts,
	}

	data, err := canonicaljson.Marshal(ik)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

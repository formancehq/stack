package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Metadata struct {
	bun.BaseModel `bun:"payments.metadata"`

	PaymentID PaymentID `bun:",pk,nullzero"`
	CreatedAt time.Time `bun:",nullzero"`
	Key       string    `bun:",pk,nullzero"`
	Value     string

	Changelog []MetadataChangelog `bun:",nullzero"`
}

type MetadataChangelog struct {
	CreatedAt time.Time `json:"createdAt"`
	Value     string    `json:"value"`
}

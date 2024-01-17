package auth

import (
	"time"

	"github.com/uptrace/bun"
)

type AccessToken struct {
	bun.BaseModel `bun:"table:access_tokens"`

	ID             string `bun:",pk"`
	ApplicationID  string
	UserID         string
	Audience       Array[string] `bun:"type:text"`
	Expiration     time.Time
	Scopes         Array[string] `bun:"type:text"`
	RefreshTokenID string        `json:"refreshTokenID"`
}

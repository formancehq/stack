package auth

import (
	"time"

	"github.com/uptrace/bun"
)

type RefreshToken struct {
	bun.BaseModel `bun:"table:refresh_tokens"`

	ID            string `bun:",pk"`
	Token         string
	AuthTime      time.Time
	AMR           Array[string] `bun:"type:text"`
	Audience      Array[string] `bun:"type:text"`
	UserID        string
	ApplicationID string
	Expiration    time.Time
	Scopes        Array[string] `bun:"type:text"`
}

type RefreshTokenRequest struct {
	RefreshToken
}

func (r *RefreshTokenRequest) GetAMR() []string {
	return r.AMR
}

func (r *RefreshTokenRequest) GetAudience() []string {
	return r.Audience
}

func (r *RefreshTokenRequest) GetAuthTime() time.Time {
	return r.AuthTime
}

func (r *RefreshTokenRequest) GetClientID() string {
	return r.ApplicationID
}

func (r *RefreshTokenRequest) GetScopes() []string {
	return r.Scopes
}

func (r *RefreshTokenRequest) GetSubject() string {
	return r.UserID
}

func (r *RefreshTokenRequest) SetCurrentScopes(scopes []string) {
	r.Scopes = scopes
}

func NewRefreshTokenRequest(r RefreshToken) *RefreshTokenRequest {
	return &RefreshTokenRequest{
		RefreshToken: r,
	}
}

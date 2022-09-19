package storage

import (
	"time"

	auth "github.com/formancehq/auth/pkg"
)

type RefreshTokenRequest struct {
	*auth.RefreshToken
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

func newRefreshTokenRequest(r *auth.RefreshToken) *RefreshTokenRequest {
	return &RefreshTokenRequest{
		RefreshToken: r,
	}
}

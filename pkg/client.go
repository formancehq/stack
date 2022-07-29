package auth

import (
	"time"

	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/oidc/pkg/op"
)

type Client struct {
	Id                             string `gorm:"primarykey"`
	Secret                         string
	RedirectURIs                   Array[string] `gorm:"type:text"`
	ApplicationType                op.ApplicationType
	AuthMethod                     oidc.AuthMethod
	ResponseTypes                  Array[oidc.ResponseType] `gorm:"type:text"`
	GrantTypes                     Array[oidc.GrantType]    `gorm:"type:text"`
	AccessTokenType                op.AccessTokenType
	DevMode                        bool
	IdTokenUserinfoClaimsAssertion bool
	ClockSkew                      time.Duration
	PostLogoutRedirectUri          Array[string] `gorm:"type:text"`
}

//WebClient will create a client of type web, which will always use Basic Auth and allow the use of refresh tokens
//user-defined redirectURIs may include:
// - http://localhost with port specification (e.g. http://localhost:9999/auth/callback)
//(the example will be used as default, if none is provided)
//func WebClient(id, secret string, redirectURIs ...string) *Client {
//	return &Client{
//		Id:                             id,
//		Secret:                         secret,
//		RedirectURIs:                   redirectURIs,
//		ApplicationType:                op.ApplicationTypeWeb,
//		AuthMethod:                     oidc.AuthMethodNone,
//		ResponseTypes:                  []oidc.ResponseType{oidc.ResponseTypeCode},
//		GrantTypes:                     []oidc.GrantType{oidc.GrantTypeCode},
//		AccessTokenType:                0,
//		DevMode:                        false, // TODO: Make configurable if required
//		IdTokenUserinfoClaimsAssertion: false,
//		ClockSkew:                      0,
//	}
//}

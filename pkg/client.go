package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/oidc/pkg/op"
)

func newHash(v string) string {
	digest := sha256.New()
	digest.Write([]byte(v))
	hash := digest.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash)
}

type ClientSecret struct {
	ID         string   `json:"id"`
	Hash       string   `json:"hash"`
	LastDigits string   `json:"lastDigits"`
	Name       string   `json:"name"`
	Metadata   Metadata `json:"metadata" gorm:"type:text"`
}

func (s ClientSecret) Check(clear string) bool {
	return s.Hash == newHash(clear)
}

func newSecret(opts SecretCreate) (ClientSecret, string) {
	clear := uuid.NewString()
	return ClientSecret{
		ID:         uuid.NewString(),
		Hash:       newHash(clear),
		LastDigits: clear[len(clear)-4:],
		Name:       opts.Name,
		Metadata:   opts.Metadata,
	}, clear
}

type Client struct {
	Id                             string              `gorm:"primarykey"`
	Secrets                        Array[ClientSecret] `gorm:"type:text"`
	RedirectURIs                   Array[string]       `gorm:"type:text"`
	ApplicationType                op.ApplicationType
	AuthMethod                     oidc.AuthMethod
	ResponseTypes                  Array[oidc.ResponseType] `gorm:"type:text"`
	GrantTypes                     Array[oidc.GrantType]    `gorm:"type:text"`
	AccessTokenType                op.AccessTokenType
	DevMode                        bool
	IdTokenUserinfoClaimsAssertion bool
	ClockSkew                      time.Duration
	PostLogoutRedirectUris         Array[string] `gorm:"type:text"`
	Scopes                         []Scope       `gorm:"many2many:client_scopes;"`
	Description                    string
	Name                           string
	Metadata                       Metadata `gorm:"type:text"`
}

func (c *Client) Update(opts ClientOptions) {
	grantTypes := []oidc.GrantType{
		oidc.GrantTypeCode,
		oidc.GrantTypeRefreshToken,
	}
	if !opts.Public {
		grantTypes = append(grantTypes, oidc.GrantTypeClientCredentials)
	}
	authMethod := oidc.AuthMethodNone
	if !opts.Public {
		authMethod = oidc.AuthMethodBasic
	}

	c.GrantTypes = grantTypes
	c.RedirectURIs = opts.RedirectUris
	c.PostLogoutRedirectUris = opts.PostLogoutRedirectUris
	c.Description = opts.Description
	c.Name = opts.Name
	c.Metadata = opts.Metadata
	c.AuthMethod = authMethod
}

func (c *Client) GenerateNewSecret(opts SecretCreate) (ClientSecret, string) {
	secret, clear := newSecret(opts)
	c.Secrets = append(c.Secrets, secret)

	return secret, clear
}

func (c *Client) HasSecret(clear string) bool {
	for _, secret := range c.Secrets {
		if secret.Check(clear) {
			return true
		}
	}
	return false
}

func (c *Client) DeleteSecret(id string) bool {
	for i, secret := range c.Secrets {
		if secret.ID == id {
			if i < len(c.Secrets)-1 {
				c.Secrets = append(c.Secrets[:i], c.Secrets[i+1:]...)
			} else {
				c.Secrets = c.Secrets[:i]
			}
			return true
		}
	}
	return false
}

func (c *Client) HasScope(id string) bool {
	for _, clientScope := range c.Scopes {
		if clientScope.ID == id {
			return true
		}
	}
	return false
}

type ClientOptions struct {
	ID                     string   `json:"id" yaml:"id"`
	Public                 bool     `json:"public" yaml:"public"`
	RedirectUris           []string `json:"redirectUris" yaml:"redirectUris"`
	Description            string   `json:"description" yaml:"description"`
	Name                   string   `json:"name" yaml:"name"`
	PostLogoutRedirectUris []string `json:"postLogoutRedirectUris" yaml:"postLogoutRedirectUris"`
	Metadata               Metadata `json:"metadata" yaml:"metadata"`
}

func NewClient(opts ClientOptions) *Client {
	if opts.ID == "" {
		opts.ID = uuid.NewString()
	}

	client := &Client{
		Id:              opts.ID,
		ApplicationType: op.ApplicationTypeWeb,
		ResponseTypes:   []oidc.ResponseType{oidc.ResponseTypeCode},
		AccessTokenType: op.AccessTokenTypeJWT,
	}
	client.Update(opts)
	return client
}

type SecretCreate struct {
	Name     string   `json:"name"`
	Metadata Metadata `json:"metadata"`
}

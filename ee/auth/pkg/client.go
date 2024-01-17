package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/uptrace/bun"

	"github.com/google/uuid"
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
	Metadata   Metadata `json:"metadata" bun:"type:text"`
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

func (c *Client) Update(opts ClientOptions) {
	c.RedirectURIs = opts.RedirectURIs
	c.PostLogoutRedirectUris = opts.PostLogoutRedirectUris
	c.Description = opts.Description
	c.Name = opts.Name
	c.Metadata = opts.Metadata
	c.Trusted = opts.Trusted
	c.Public = opts.Public
	c.Scopes = opts.Scopes
}

func (c *Client) GenerateNewSecret(opts SecretCreate) (ClientSecret, string) {
	secret, clear := newSecret(opts)
	c.Secrets = append(c.Secrets, secret)

	return secret, clear
}

func (c *Client) ValidateSecret(secret string) error {
	if !c.HasSecret(secret) {
		return errors.New("invalid secret")
	}
	return nil
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
		if clientScope == id {
			return true
		}
	}
	return false
}

type Client struct {
	bun.BaseModel `bun:"table:clients"`

	ClientOptions
	Secrets Array[ClientSecret] `bun:",type:text" json:"secrets"`
}

func (c *Client) GetScopes() []string {
	return c.Scopes
}

type StaticClient struct {
	ClientOptions `mapstructure:",squash" yaml:",inline"`
	Secrets       []string `json:"secrets" yaml:"secrets"`
}

func (s StaticClient) ValidateSecret(secret string) error {
	for _, clientSecret := range s.Secrets {
		if clientSecret == secret {
			return nil
		}
	}
	return errors.New("invalid secret")
}

func (s StaticClient) GetScopes() []string {
	return s.Scopes
}

type ClientOptions struct {
	Id                     string        `json:"id" yaml:"id"`
	Public                 bool          `json:"public" yaml:"public"`
	RedirectURIs           Array[string] `json:"redirectUris" yaml:"redirectUris" bun:"redirect_uris,type:text"`
	Description            string        `json:"description" yaml:"description"`
	Name                   string        `json:"name" yaml:"name"`
	PostLogoutRedirectUris Array[string] `json:"postLogoutRedirectUris" yaml:"postLogoutRedirectUris" bun:"post_logout_redirect_uris,type:text"`
	Metadata               Metadata      `json:"metadata" yaml:"metadata" bun:"type:text"`
	Trusted                bool          `json:"trusted" yaml:"trusted"`
	Scopes                 Array[string] `bun:"type:text" json:"scopes" yaml:"scopes"`
}

func (s *ClientOptions) IsTrusted() bool {
	return s.Trusted
}

func (c *ClientOptions) GetID() string {
	return c.Id
}

func (c *ClientOptions) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *ClientOptions) GetPostLogoutRedirectUris() []string {
	return c.PostLogoutRedirectUris
}

func (c *ClientOptions) IsPublic() bool {
	return c.Public
}

func NewClient(opts ClientOptions) *Client {
	if opts.Id == "" {
		opts.Id = uuid.NewString()
	}

	client := &Client{
		ClientOptions: opts,
	}
	client.Update(opts)
	return client
}

type SecretCreate struct {
	Name     string   `json:"name"`
	Metadata Metadata `json:"metadata"`
}

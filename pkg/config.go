package webhooks

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type ConfigUser struct {
	Endpoint   string   `json:"endpoint" bson:"endpoint"`
	Secret     string   `json:"secret" bson:"secret"`
	EventTypes []string `json:"eventTypes" bson:"eventTypes"`
}

type Config struct {
	ConfigUser `bson:"inline"`
	ID         string    `json:"_id" bson:"_id"`
	Active     bool      `json:"active" bson:"active"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
}

const (
	KeySecret     = "secret"
	KeyEventTypes = "eventTypes"
	KeyID         = "_id"
	KeyActive     = "active"
	KeyUpdatedAt  = "updatedAt"
)

var (
	ErrInvalidEndpoint   = errors.New("endpoint should be a valid url")
	ErrInvalidEventTypes = errors.New("eventTypes should be filled")
	ErrInvalidSecret     = errors.New("decoded secret should be of size 24")
)

func (c *ConfigUser) Validate() error {
	if u, err := url.Parse(c.Endpoint); err != nil || len(u.String()) == 0 {
		return ErrInvalidEndpoint
	}

	if c.Secret == "" {
		c.Secret = NewSecret()
	} else {
		if decoded, err := base64.StdEncoding.DecodeString(c.Secret); err != nil {
			return fmt.Errorf("secret should be base64 encoded: %w", err)
		} else if len(decoded) != 24 {
			return ErrInvalidSecret
		}
	}

	if len(c.EventTypes) == 0 {
		return ErrInvalidEventTypes
	}

	for i, t := range c.EventTypes {
		if len(t) == 0 {
			return ErrInvalidEventTypes
		}
		c.EventTypes[i] = strings.ToLower(t)
	}

	return nil
}

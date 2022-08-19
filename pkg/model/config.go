package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	Endpoint   string   `json:"endpoint" bson:"endpoint"`
	Secret     string   `json:"secret" bson:"secret"`
	EventTypes []string `json:"eventTypes" bson:"eventTypes"`
}

type ConfigInserted struct {
	Config    `bson:"inline"`
	ID        string    `json:"_id" bson:"_id"`
	Active    bool      `json:"active" bson:"active"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (c *Config) Validate() error {
	if u, err := url.Parse(c.Endpoint); err != nil || len(u.String()) == 0 {
		return errors.New("endpoint should be a valid url")
	}

	if c.Secret == "" {
		c.Secret = NewSecret()
	} else {
		var decoded []byte
		var err error
		if decoded, err = base64.StdEncoding.DecodeString(c.Secret); err != nil {
			return fmt.Errorf("secret should be base64 encoded: %w", err)
		}
		if len(decoded) != 24 {
			return fmt.Errorf("decoded secret should be of size 24: actual size %d", len(decoded))
		}
	}

	if len(c.EventTypes) == 0 {
		return errors.New("eventTypes should be filled")
	}

	for _, t := range c.EventTypes {
		if len(t) == 0 {
			return errors.New("eventTypes should be filled")
		}
	}

	return nil
}

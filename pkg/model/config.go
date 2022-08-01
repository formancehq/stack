package model

import (
	"errors"
	"net/url"
	"time"
)

type Config struct {
	Active     bool     `json:"active" bson:"active"`
	EventTypes []string `json:"eventTypes,omitempty" bson:"eventTypes,omitempty"`
	Endpoints  []string `json:"endpoints,omitempty" bson:"endpoints,omitempty"`
}

type ConfigInserted struct {
	Config     `bson:"inline"`
	ID         string    `json:"_id" bson:"_id"`
	InsertedAt time.Time `json:"insertedAt" bson:"insertedAt"`
}

func (c Config) Validate() error {
	if c.Active {
		if len(c.EventTypes) < 1 || len(c.Endpoints) < 1 {
			return errors.New(
				"the body should have at least one type of events and one endpoint")
		}

		for _, endpoint := range c.Endpoints {
			if _, err := url.Parse(endpoint); err != nil {
				return errors.New(
					"endpoints should be valid urls")
			}
		}
	} else {
		if len(c.EventTypes) > 0 || len(c.Endpoints) > 0 {
			return errors.New(
				"the body to set a webhook inactive shouldn't contain any types of events or endpoints")
		}
	}

	return nil
}

package fawry

import "encoding/json"

type Config struct {
	Name         string `json:"name" yaml:"name" bson:"name"`
	Endpoint     string `json:"endpoint" yaml:"endpoint" bson:"endpoint"`
	SecureKey    string `json:"secureKey" yaml:"secureKey" bson:"secureKey"`
	MerchantCode string `json:"merchantCode" yaml:"merchantCode" bson:"merchantCode"`
}

func (c Config) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c Config) ConnectorName() string {
	return c.Name
}

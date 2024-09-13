package plugins

import (
	_ "embed"
	"encoding/json"
	"errors"
)

//go:embed configs.json
var configsFile []byte

type Type string

const (
	TypeLongString              Type = "long string"
	TypeString                  Type = "string"
	TypeDurationNs              Type = "duration ns"
	TypeDurationUnsignedInteger Type = "unsigned integer"
	TypeBoolean                 Type = "boolean"
)

type Configs map[string]Config
type Config map[string]Parameter
type Parameter struct {
	DataType     Type   `json:"dataType"`
	Required     bool   `json:"required"`
	DefaultValue string `json:"defaultValue"`
}

var (
	defaultParameters = map[string]Parameter{
		"pollingPeriod": {
			DataType:     "duration ns",
			Required:     false,
			DefaultValue: "2m",
		},
		"pageSize": {
			DataType:     "unsigned integer",
			Required:     false,
			DefaultValue: "100",
		},
		"name": {
			DataType: "string",
			Required: true,
		},
	}

	configs Configs
)

func init() {
	if err := json.Unmarshal(configsFile, &configs); err != nil {
		panic(err)
	}

	for key := range configs {
		for paramName, param := range defaultParameters {
			if _, ok := configs[key][paramName]; !ok {
				configs[key][paramName] = param
			}
		}
	}
}

func GetConfigs() Configs {
	return configs
}

func GetConfig(provider string) (Config, error) {
	config, ok := configs[provider]
	if !ok {
		return nil, errors.New("config not found")
	}

	return config, nil
}

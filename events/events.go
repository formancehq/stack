package events

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v3"
)

//go:embed v1
var V1 embed.FS

func ComputeSchema(serviceName, eventName string) (*gojsonschema.Schema, error) {
	baseData, err := fs.ReadFile(V1, "v1/base/base.yaml")
	if err != nil {
		return nil, err
	}
	base := map[string]any{}
	if err := yaml.Unmarshal(baseData, &base); err != nil {
		return nil, err
	}

	eventData, err := fs.ReadFile(V1, fmt.Sprintf("v1/%s/%s.yaml", serviceName, eventName))
	if err != nil {
		return nil, err
	}
	event := map[string]any{}
	if err := yaml.Unmarshal(eventData, &event); err != nil {
		return nil, err
	}

	base["properties"].(map[string]any)["payload"] = event

	loader := gojsonschema.NewGoLoader(base)
	return gojsonschema.NewSchema(loader)
}

func Check(data []byte, serviceName, eventName string) error {
	schema, err := ComputeSchema(serviceName, eventName)
	if err != nil {
		return errors.Wrap(err, "computing schema")
	}
	result, err := schema.Validate(gojsonschema.NewStringLoader(string(data)))
	if err != nil {
		return errors.Wrap(err, "validating schema")
	}
	if len(result.Errors()) > 0 {
		ret := ""
		for _, resultError := range result.Errors() {
			ret += resultError.String() + "\r\n"
		}
		return errors.New(ret)
	}
	return nil
}

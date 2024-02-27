package events

import (
	"fmt"
	"path/filepath"

	"embed"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
)

//go:embed base.yaml
var baseEvent string

//go:embed services
var services embed.FS

func ComputeSchema(serviceName, eventName string) (*gojsonschema.Schema, error) {
	base := map[string]any{}
	if err := yaml.Unmarshal([]byte(baseEvent), &base); err != nil {
		return nil, err
	}

	ls, err := services.ReadDir(filepath.Join("services", serviceName))
	if err != nil {
		return nil, errors.Wrapf(err, "reading events directory for service '%s'", serviceName)
	}

	var moreRecent string
	for _, directory := range ls {
		if moreRecent == "" || semver.Compare(directory.Name(), moreRecent) > 0 {
			moreRecent = directory.Name()
		}
	}

	if moreRecent == "" {
		return nil, fmt.Errorf("error retrieving more recent version directory for service '%s'", serviceName)
	}

	eventData, err := services.ReadFile(fmt.Sprintf("services/%s/%s/%s.yaml", serviceName, moreRecent, eventName))
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

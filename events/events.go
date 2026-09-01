package events

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"

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
	ls, err := services.ReadDir(filepath.Join("services", serviceName))
	if err != nil {
		return nil, errors.Wrapf(err, "reading events directory for service '%s'", serviceName)
	}

	versions := make([]string, 0, len(ls))
	for _, directory := range ls {
		if directory.IsDir() && semver.IsValid(directory.Name()) {
			versions = append(versions, directory.Name())
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return semver.Compare(versions[i], versions[j]) > 0
	})

	if len(versions) == 0 {
		return nil, fmt.Errorf("error retrieving more recent version directory for service '%s'", serviceName)
	}

	for _, version := range versions {
		_, err := services.ReadFile(fmt.Sprintf("services/%s/%s/%s.yaml", serviceName, version, eventName))
		if err == nil {
			return ComputeSchemaForVersion(serviceName, version, eventName)
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("event schema '%s' not found for service '%s'", eventName, serviceName)
}

// ComputeSchemaForVersion returns the schema for an exact service version.
// Versions with a base.yaml describe a complete event envelope and compose
// the event-specific constraints through allOf. Older versions keep using the
// historical shared envelope with the event schema injected under payload.
func ComputeSchemaForVersion(serviceName, version, eventName string) (*gojsonschema.Schema, error) {
	eventData, err := services.ReadFile(fmt.Sprintf("services/%s/%s/%s.yaml", serviceName, version, eventName))
	if err != nil {
		return nil, err
	}

	event := map[string]any{}
	if err := yaml.Unmarshal(eventData, &event); err != nil {
		return nil, err
	}

	versionBaseData, err := services.ReadFile(fmt.Sprintf("services/%s/%s/base.yaml", serviceName, version))
	if err == nil {
		base := map[string]any{}
		if err := yaml.Unmarshal(versionBaseData, &base); err != nil {
			return nil, err
		}

		allOf, _ := base["allOf"].([]any)
		base["allOf"] = append(allOf, event)

		return compileSchema(base)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	base := map[string]any{}
	if err := yaml.Unmarshal([]byte(baseEvent), &base); err != nil {
		return nil, err
	}

	properties, ok := base["properties"].(map[string]any)
	if !ok {
		return nil, errors.New("base event schema has no properties object")
	}
	properties["payload"] = event

	return compileSchema(base)
}

func compileSchema(schema map[string]any) (*gojsonschema.Schema, error) {
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling schema")
	}
	return gojsonschema.NewSchema(gojsonschema.NewBytesLoader(data))
}

func Check(data []byte, serviceName, eventName string) error {
	schema, err := ComputeSchema(serviceName, eventName)
	if err != nil {
		return errors.Wrap(err, "computing schema")
	}

	return validate(data, schema)
}

// CheckForVersion validates an event against an exact service version.
func CheckForVersion(data []byte, serviceName, version, eventName string) error {
	schema, err := ComputeSchemaForVersion(serviceName, version, eventName)
	if err != nil {
		return errors.Wrap(err, "computing schema")
	}

	return validate(data, schema)
}

func validate(data []byte, schema *gojsonschema.Schema) error {
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

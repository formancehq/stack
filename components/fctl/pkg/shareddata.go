package fctl

import (
	"encoding/json"

	"github.com/TylerBrock/colorjson"
	"github.com/pkg/errors"
)

type SharedStore struct {
	data    interface{}
	profile *Profile
	config  *Config

	// Those data are not printed in the json output
	additionnalData map[string]interface{}
	//additionnalKeyType 	map[string]

}

var sharedStore = &SharedStore{
	additionnalData: make(map[string]interface{}),
}

// GetSharedData returns the shared data store
func GetSharedData() interface{} {
	return sharedStore.data
}

func GetSharedProfile() *Profile {
	return sharedStore.profile
}

func GetSharedConfig() *Config {
	return sharedStore.config
}

func SetSharedData(data interface{}, profile *Profile, config *Config, additionnalData map[string]interface{}) {
	sharedStore.data = data
	sharedStore.profile = profile
	sharedStore.config = config
	sharedStore.additionnalData = additionnalData
}

func SetSharedAdditionnalData(key string, value interface{}) {
	sharedStore.additionnalData[key] = value
}

func GetSharedAdditionnalData(key string) interface{} {
	return sharedStore.additionnalData[key]
}

type ExportedData struct {
	Data interface{} `json:"data"`
}

func ShareStoreToJson() ([]byte, error) {
	if (sharedStore.data) == nil {
		return nil, errors.New("no data to marshal")
	}

	// Inject into export struct
	export := ExportedData{
		Data: sharedStore.data,
	}

	// Marshal to JSON then print to stdout
	s, err := json.Marshal(export)
	if err != nil {
		return nil, err
	}

	raw := make(map[string]any)
	if err := json.Unmarshal(s, &raw); err == nil {
		f := colorjson.NewFormatter()
		f.Indent = 2
		colorized, err := f.Marshal(raw)
		if err != nil {
			panic(err)
		}
		return colorized, nil
	} else {
		return s, nil
	}

}

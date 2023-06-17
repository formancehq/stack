package fctl

import (
	"encoding/json"

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

func ShareStoreToJson() ([]byte, error) {
	if (sharedStore.data) == nil {
		errors.New("no data to marshal")
	}

	// Marshal to JSON then print to stdout
	return json.Marshal(sharedStore.data)
}

package fctl

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type SharedStore struct {
	data    interface{}
	profile *Profile
	config  *Config
}

var sharedStore = &SharedStore{}

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

func SetSharedData(data interface{}, profile *Profile, config *Config) {
	sharedStore.data = data
	sharedStore.profile = profile
	sharedStore.config = config
}

func ShareStoreToJson() ([]byte, error) {
	if (sharedStore.data) == nil {
		errors.New("no data to marshal")
	}

	// Marshal to JSON then print to stdout
	return json.Marshal(sharedStore.data)
}

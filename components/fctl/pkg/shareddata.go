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

func NewSharedStore() *SharedStore {
	return &SharedStore{
		additionnalData: make(map[string]interface{}),
	}
}

// GetSharedData returns the shared data store
func (s *SharedStore) GetData() interface{} {
	return s.data
}

func (s *SharedStore) GetProfile() *Profile {
	return s.profile
}

func (s *SharedStore) GetConfig() *Config {
	return s.config
}

func (s *SharedStore) SetConfig(c *Config) *SharedStore {
	s.config = c
	return s
}

func (s *SharedStore) SetData(data interface{}) *SharedStore {
	s.data = data
	return s
}

func (s *SharedStore) SetProfile(p *Profile) *SharedStore {
	s.profile = p
	return s
}

func (s *SharedStore) SetAdditionnalData(key string, value interface{}) {
	s.additionnalData[key] = value
}

func (s *SharedStore) GetAdditionnalData(key string) interface{} {
	return s.additionnalData[key]
}

type ExportedData struct {
	Data interface{} `json:"data"`
}

func (s *SharedStore) ToJson() ([]byte, error) {
	if (s.data) == nil {
		return nil, errors.New("no data to marshal")
	}

	// Inject into export struct
	export := ExportedData{
		Data: s.data,
	}

	// Marshal to JSON then print to stdout
	out, err := json.Marshal(export)
	if err != nil {
		return nil, err
	}

	raw := make(map[string]any)
	if err := json.Unmarshal(out, &raw); err == nil {
		f := colorjson.NewFormatter()
		f.Indent = 2
		colorized, err := f.Marshal(raw)
		if err != nil {
			panic(err)
		}
		return colorized, nil
	} else {
		return out, nil
	}

}

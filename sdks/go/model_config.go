/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: v1.0.20230301
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
)

// checks if the Config type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Config{}

// Config struct for Config
type Config struct {
	Storage LedgerStorage `json:"storage"`
}

// NewConfig instantiates a new Config object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfig(storage LedgerStorage) *Config {
	this := Config{}
	this.Storage = storage
	return &this
}

// NewConfigWithDefaults instantiates a new Config object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigWithDefaults() *Config {
	this := Config{}
	return &this
}

// GetStorage returns the Storage field value
func (o *Config) GetStorage() LedgerStorage {
	if o == nil {
		var ret LedgerStorage
		return ret
	}

	return o.Storage
}

// GetStorageOk returns a tuple with the Storage field value
// and a boolean to check if the value has been set.
func (o *Config) GetStorageOk() (*LedgerStorage, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Storage, true
}

// SetStorage sets field value
func (o *Config) SetStorage(v LedgerStorage) {
	o.Storage = v
}

func (o Config) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Config) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["storage"] = o.Storage
	return toSerialize, nil
}

type NullableConfig struct {
	value *Config
	isSet bool
}

func (v NullableConfig) Get() *Config {
	return v.value
}

func (v *NullableConfig) Set(val *Config) {
	v.value = val
	v.isSet = true
}

func (v NullableConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfig(val *Config) *NullableConfig {
	return &NullableConfig{value: val, isSet: true}
}

func (v NullableConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



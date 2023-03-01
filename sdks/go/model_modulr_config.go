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

// checks if the ModulrConfig type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ModulrConfig{}

// ModulrConfig struct for ModulrConfig
type ModulrConfig struct {
	ApiKey string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
	Endpoint *string `json:"endpoint,omitempty"`
}

// NewModulrConfig instantiates a new ModulrConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewModulrConfig(apiKey string, apiSecret string) *ModulrConfig {
	this := ModulrConfig{}
	this.ApiKey = apiKey
	this.ApiSecret = apiSecret
	return &this
}

// NewModulrConfigWithDefaults instantiates a new ModulrConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewModulrConfigWithDefaults() *ModulrConfig {
	this := ModulrConfig{}
	return &this
}

// GetApiKey returns the ApiKey field value
func (o *ModulrConfig) GetApiKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ApiKey
}

// GetApiKeyOk returns a tuple with the ApiKey field value
// and a boolean to check if the value has been set.
func (o *ModulrConfig) GetApiKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ApiKey, true
}

// SetApiKey sets field value
func (o *ModulrConfig) SetApiKey(v string) {
	o.ApiKey = v
}

// GetApiSecret returns the ApiSecret field value
func (o *ModulrConfig) GetApiSecret() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ApiSecret
}

// GetApiSecretOk returns a tuple with the ApiSecret field value
// and a boolean to check if the value has been set.
func (o *ModulrConfig) GetApiSecretOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ApiSecret, true
}

// SetApiSecret sets field value
func (o *ModulrConfig) SetApiSecret(v string) {
	o.ApiSecret = v
}

// GetEndpoint returns the Endpoint field value if set, zero value otherwise.
func (o *ModulrConfig) GetEndpoint() string {
	if o == nil || IsNil(o.Endpoint) {
		var ret string
		return ret
	}
	return *o.Endpoint
}

// GetEndpointOk returns a tuple with the Endpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ModulrConfig) GetEndpointOk() (*string, bool) {
	if o == nil || IsNil(o.Endpoint) {
		return nil, false
	}
	return o.Endpoint, true
}

// HasEndpoint returns a boolean if a field has been set.
func (o *ModulrConfig) HasEndpoint() bool {
	if o != nil && !IsNil(o.Endpoint) {
		return true
	}

	return false
}

// SetEndpoint gets a reference to the given string and assigns it to the Endpoint field.
func (o *ModulrConfig) SetEndpoint(v string) {
	o.Endpoint = &v
}

func (o ModulrConfig) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ModulrConfig) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["apiKey"] = o.ApiKey
	toSerialize["apiSecret"] = o.ApiSecret
	if !IsNil(o.Endpoint) {
		toSerialize["endpoint"] = o.Endpoint
	}
	return toSerialize, nil
}

type NullableModulrConfig struct {
	value *ModulrConfig
	isSet bool
}

func (v NullableModulrConfig) Get() *ModulrConfig {
	return v.value
}

func (v *NullableModulrConfig) Set(val *ModulrConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableModulrConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableModulrConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableModulrConfig(val *ModulrConfig) *NullableModulrConfig {
	return &NullableModulrConfig{value: val, isSet: true}
}

func (v NullableModulrConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableModulrConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



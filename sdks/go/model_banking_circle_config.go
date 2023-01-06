/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions />

API version: develop
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
)

// BankingCircleConfig struct for BankingCircleConfig
type BankingCircleConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Endpoint string `json:"endpoint"`
	AuthorizationEndpoint string `json:"authorizationEndpoint"`
}

// NewBankingCircleConfig instantiates a new BankingCircleConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBankingCircleConfig(username string, password string, endpoint string, authorizationEndpoint string) *BankingCircleConfig {
	this := BankingCircleConfig{}
	this.Username = username
	this.Password = password
	this.Endpoint = endpoint
	this.AuthorizationEndpoint = authorizationEndpoint
	return &this
}

// NewBankingCircleConfigWithDefaults instantiates a new BankingCircleConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBankingCircleConfigWithDefaults() *BankingCircleConfig {
	this := BankingCircleConfig{}
	return &this
}

// GetUsername returns the Username field value
func (o *BankingCircleConfig) GetUsername() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Username
}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
func (o *BankingCircleConfig) GetUsernameOk() (*string, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Username, true
}

// SetUsername sets field value
func (o *BankingCircleConfig) SetUsername(v string) {
	o.Username = v
}

// GetPassword returns the Password field value
func (o *BankingCircleConfig) GetPassword() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Password
}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
func (o *BankingCircleConfig) GetPasswordOk() (*string, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Password, true
}

// SetPassword sets field value
func (o *BankingCircleConfig) SetPassword(v string) {
	o.Password = v
}

// GetEndpoint returns the Endpoint field value
func (o *BankingCircleConfig) GetEndpoint() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Endpoint
}

// GetEndpointOk returns a tuple with the Endpoint field value
// and a boolean to check if the value has been set.
func (o *BankingCircleConfig) GetEndpointOk() (*string, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Endpoint, true
}

// SetEndpoint sets field value
func (o *BankingCircleConfig) SetEndpoint(v string) {
	o.Endpoint = v
}

// GetAuthorizationEndpoint returns the AuthorizationEndpoint field value
func (o *BankingCircleConfig) GetAuthorizationEndpoint() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.AuthorizationEndpoint
}

// GetAuthorizationEndpointOk returns a tuple with the AuthorizationEndpoint field value
// and a boolean to check if the value has been set.
func (o *BankingCircleConfig) GetAuthorizationEndpointOk() (*string, bool) {
	if o == nil {
    return nil, false
	}
	return &o.AuthorizationEndpoint, true
}

// SetAuthorizationEndpoint sets field value
func (o *BankingCircleConfig) SetAuthorizationEndpoint(v string) {
	o.AuthorizationEndpoint = v
}

func (o BankingCircleConfig) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["username"] = o.Username
	}
	if true {
		toSerialize["password"] = o.Password
	}
	if true {
		toSerialize["endpoint"] = o.Endpoint
	}
	if true {
		toSerialize["authorizationEndpoint"] = o.AuthorizationEndpoint
	}
	return json.Marshal(toSerialize)
}

type NullableBankingCircleConfig struct {
	value *BankingCircleConfig
	isSet bool
}

func (v NullableBankingCircleConfig) Get() *BankingCircleConfig {
	return v.value
}

func (v *NullableBankingCircleConfig) Set(val *BankingCircleConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableBankingCircleConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableBankingCircleConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBankingCircleConfig(val *BankingCircleConfig) *NullableBankingCircleConfig {
	return &NullableBankingCircleConfig{value: val, isSet: true}
}

func (v NullableBankingCircleConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBankingCircleConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

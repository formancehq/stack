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

// checks if the WalletSource type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WalletSource{}

// WalletSource struct for WalletSource
type WalletSource struct {
	Balance string `json:"balance"`
	Id string `json:"id"`
}

// NewWalletSource instantiates a new WalletSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWalletSource(balance string, id string) *WalletSource {
	this := WalletSource{}
	this.Balance = balance
	this.Id = id
	return &this
}

// NewWalletSourceWithDefaults instantiates a new WalletSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWalletSourceWithDefaults() *WalletSource {
	this := WalletSource{}
	return &this
}

// GetBalance returns the Balance field value
func (o *WalletSource) GetBalance() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Balance
}

// GetBalanceOk returns a tuple with the Balance field value
// and a boolean to check if the value has been set.
func (o *WalletSource) GetBalanceOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Balance, true
}

// SetBalance sets field value
func (o *WalletSource) SetBalance(v string) {
	o.Balance = v
}

// GetId returns the Id field value
func (o *WalletSource) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WalletSource) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WalletSource) SetId(v string) {
	o.Id = v
}

func (o WalletSource) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WalletSource) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["balance"] = o.Balance
	toSerialize["id"] = o.Id
	return toSerialize, nil
}

type NullableWalletSource struct {
	value *WalletSource
	isSet bool
}

func (v NullableWalletSource) Get() *WalletSource {
	return v.value
}

func (v *NullableWalletSource) Set(val *WalletSource) {
	v.value = val
	v.isSet = true
}

func (v NullableWalletSource) IsSet() bool {
	return v.isSet
}

func (v *NullableWalletSource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWalletSource(val *WalletSource) *NullableWalletSource {
	return &NullableWalletSource{value: val, isSet: true}
}

func (v NullableWalletSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWalletSource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



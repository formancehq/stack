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

// ListWalletsResponseCursorAllOf struct for ListWalletsResponseCursorAllOf
type ListWalletsResponseCursorAllOf struct {
	Data []Wallet `json:"data"`
}

// NewListWalletsResponseCursorAllOf instantiates a new ListWalletsResponseCursorAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListWalletsResponseCursorAllOf(data []Wallet) *ListWalletsResponseCursorAllOf {
	this := ListWalletsResponseCursorAllOf{}
	this.Data = data
	return &this
}

// NewListWalletsResponseCursorAllOfWithDefaults instantiates a new ListWalletsResponseCursorAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListWalletsResponseCursorAllOfWithDefaults() *ListWalletsResponseCursorAllOf {
	this := ListWalletsResponseCursorAllOf{}
	return &this
}

// GetData returns the Data field value
func (o *ListWalletsResponseCursorAllOf) GetData() []Wallet {
	if o == nil {
		var ret []Wallet
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *ListWalletsResponseCursorAllOf) GetDataOk() ([]Wallet, bool) {
	if o == nil {
    return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *ListWalletsResponseCursorAllOf) SetData(v []Wallet) {
	o.Data = v
}

func (o ListWalletsResponseCursorAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableListWalletsResponseCursorAllOf struct {
	value *ListWalletsResponseCursorAllOf
	isSet bool
}

func (v NullableListWalletsResponseCursorAllOf) Get() *ListWalletsResponseCursorAllOf {
	return v.value
}

func (v *NullableListWalletsResponseCursorAllOf) Set(val *ListWalletsResponseCursorAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableListWalletsResponseCursorAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableListWalletsResponseCursorAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListWalletsResponseCursorAllOf(val *ListWalletsResponseCursorAllOf) *NullableListWalletsResponseCursorAllOf {
	return &NullableListWalletsResponseCursorAllOf{value: val, isSet: true}
}

func (v NullableListWalletsResponseCursorAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListWalletsResponseCursorAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

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

// GetBalanceResponse struct for GetBalanceResponse
type GetBalanceResponse struct {
	Data BalanceWithAssets `json:"data"`
}

// NewGetBalanceResponse instantiates a new GetBalanceResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetBalanceResponse(data BalanceWithAssets) *GetBalanceResponse {
	this := GetBalanceResponse{}
	this.Data = data
	return &this
}

// NewGetBalanceResponseWithDefaults instantiates a new GetBalanceResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetBalanceResponseWithDefaults() *GetBalanceResponse {
	this := GetBalanceResponse{}
	return &this
}

// GetData returns the Data field value
func (o *GetBalanceResponse) GetData() BalanceWithAssets {
	if o == nil {
		var ret BalanceWithAssets
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *GetBalanceResponse) GetDataOk() (*BalanceWithAssets, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *GetBalanceResponse) SetData(v BalanceWithAssets) {
	o.Data = v
}

func (o GetBalanceResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableGetBalanceResponse struct {
	value *GetBalanceResponse
	isSet bool
}

func (v NullableGetBalanceResponse) Get() *GetBalanceResponse {
	return v.value
}

func (v *NullableGetBalanceResponse) Set(val *GetBalanceResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetBalanceResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetBalanceResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetBalanceResponse(val *GetBalanceResponse) *NullableGetBalanceResponse {
	return &NullableGetBalanceResponse{value: val, isSet: true}
}

func (v NullableGetBalanceResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetBalanceResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

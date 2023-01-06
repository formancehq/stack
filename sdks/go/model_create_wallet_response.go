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

// CreateWalletResponse struct for CreateWalletResponse
type CreateWalletResponse struct {
	Data Wallet `json:"data"`
}

// NewCreateWalletResponse instantiates a new CreateWalletResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateWalletResponse(data Wallet) *CreateWalletResponse {
	this := CreateWalletResponse{}
	this.Data = data
	return &this
}

// NewCreateWalletResponseWithDefaults instantiates a new CreateWalletResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateWalletResponseWithDefaults() *CreateWalletResponse {
	this := CreateWalletResponse{}
	return &this
}

// GetData returns the Data field value
func (o *CreateWalletResponse) GetData() Wallet {
	if o == nil {
		var ret Wallet
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *CreateWalletResponse) GetDataOk() (*Wallet, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *CreateWalletResponse) SetData(v Wallet) {
	o.Data = v
}

func (o CreateWalletResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableCreateWalletResponse struct {
	value *CreateWalletResponse
	isSet bool
}

func (v NullableCreateWalletResponse) Get() *CreateWalletResponse {
	return v.value
}

func (v *NullableCreateWalletResponse) Set(val *CreateWalletResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateWalletResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateWalletResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateWalletResponse(val *CreateWalletResponse) *NullableCreateWalletResponse {
	return &NullableCreateWalletResponse{value: val, isSet: true}
}

func (v NullableCreateWalletResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateWalletResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

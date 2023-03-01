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

// checks if the GetWalletResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetWalletResponse{}

// GetWalletResponse struct for GetWalletResponse
type GetWalletResponse struct {
	Data WalletWithBalances `json:"data"`
}

// NewGetWalletResponse instantiates a new GetWalletResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetWalletResponse(data WalletWithBalances) *GetWalletResponse {
	this := GetWalletResponse{}
	this.Data = data
	return &this
}

// NewGetWalletResponseWithDefaults instantiates a new GetWalletResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetWalletResponseWithDefaults() *GetWalletResponse {
	this := GetWalletResponse{}
	return &this
}

// GetData returns the Data field value
func (o *GetWalletResponse) GetData() WalletWithBalances {
	if o == nil {
		var ret WalletWithBalances
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *GetWalletResponse) GetDataOk() (*WalletWithBalances, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *GetWalletResponse) SetData(v WalletWithBalances) {
	o.Data = v
}

func (o GetWalletResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetWalletResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableGetWalletResponse struct {
	value *GetWalletResponse
	isSet bool
}

func (v NullableGetWalletResponse) Get() *GetWalletResponse {
	return v.value
}

func (v *NullableGetWalletResponse) Set(val *GetWalletResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetWalletResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetWalletResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetWalletResponse(val *GetWalletResponse) *NullableGetWalletResponse {
	return &NullableGetWalletResponse{value: val, isSet: true}
}

func (v NullableGetWalletResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetWalletResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



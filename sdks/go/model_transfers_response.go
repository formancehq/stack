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

// checks if the TransfersResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TransfersResponse{}

// TransfersResponse struct for TransfersResponse
type TransfersResponse struct {
	Data []TransfersResponseDataInner `json:"data,omitempty"`
}

// NewTransfersResponse instantiates a new TransfersResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTransfersResponse() *TransfersResponse {
	this := TransfersResponse{}
	return &this
}

// NewTransfersResponseWithDefaults instantiates a new TransfersResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTransfersResponseWithDefaults() *TransfersResponse {
	this := TransfersResponse{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *TransfersResponse) GetData() []TransfersResponseDataInner {
	if o == nil || IsNil(o.Data) {
		var ret []TransfersResponseDataInner
		return ret
	}
	return o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TransfersResponse) GetDataOk() ([]TransfersResponseDataInner, bool) {
	if o == nil || IsNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *TransfersResponse) HasData() bool {
	if o != nil && !IsNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given []TransfersResponseDataInner and assigns it to the Data field.
func (o *TransfersResponse) SetData(v []TransfersResponseDataInner) {
	o.Data = v
}

func (o TransfersResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TransfersResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return toSerialize, nil
}

type NullableTransfersResponse struct {
	value *TransfersResponse
	isSet bool
}

func (v NullableTransfersResponse) Get() *TransfersResponse {
	return v.value
}

func (v *NullableTransfersResponse) Set(val *TransfersResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTransfersResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTransfersResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTransfersResponse(val *TransfersResponse) *NullableTransfersResponse {
	return &NullableTransfersResponse{value: val, isSet: true}
}

func (v NullableTransfersResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTransfersResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



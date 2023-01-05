/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package membershipclient

import (
	"encoding/json"
)

// checks if the ListStacksResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ListStacksResponse{}

// ListStacksResponse struct for ListStacksResponse
type ListStacksResponse struct {
	Data []Stack `json:"data,omitempty"`
}

// NewListStacksResponse instantiates a new ListStacksResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListStacksResponse() *ListStacksResponse {
	this := ListStacksResponse{}
	return &this
}

// NewListStacksResponseWithDefaults instantiates a new ListStacksResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListStacksResponseWithDefaults() *ListStacksResponse {
	this := ListStacksResponse{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *ListStacksResponse) GetData() []Stack {
	if o == nil || isNil(o.Data) {
		var ret []Stack
		return ret
	}
	return o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ListStacksResponse) GetDataOk() ([]Stack, bool) {
	if o == nil || isNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *ListStacksResponse) HasData() bool {
	if o != nil && !isNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given []Stack and assigns it to the Data field.
func (o *ListStacksResponse) SetData(v []Stack) {
	o.Data = v
}

func (o ListStacksResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ListStacksResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return toSerialize, nil
}

type NullableListStacksResponse struct {
	value *ListStacksResponse
	isSet bool
}

func (v NullableListStacksResponse) Get() *ListStacksResponse {
	return v.value
}

func (v *NullableListStacksResponse) Set(val *ListStacksResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableListStacksResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableListStacksResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListStacksResponse(val *ListStacksResponse) *NullableListStacksResponse {
	return &NullableListStacksResponse{value: val, isSet: true}
}

func (v NullableListStacksResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListStacksResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



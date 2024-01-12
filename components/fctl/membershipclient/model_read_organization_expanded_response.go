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

// checks if the ReadOrganizationExpandedResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ReadOrganizationExpandedResponse{}

// ReadOrganizationExpandedResponse struct for ReadOrganizationExpandedResponse
type ReadOrganizationExpandedResponse struct {
	Data *OrganizationExpanded `json:"data,omitempty"`
}

// NewReadOrganizationExpandedResponse instantiates a new ReadOrganizationExpandedResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReadOrganizationExpandedResponse() *ReadOrganizationExpandedResponse {
	this := ReadOrganizationExpandedResponse{}
	return &this
}

// NewReadOrganizationExpandedResponseWithDefaults instantiates a new ReadOrganizationExpandedResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReadOrganizationExpandedResponseWithDefaults() *ReadOrganizationExpandedResponse {
	this := ReadOrganizationExpandedResponse{}
	return &this
}

// GetData returns the Data field value if set, zero value otherwise.
func (o *ReadOrganizationExpandedResponse) GetData() OrganizationExpanded {
	if o == nil || IsNil(o.Data) {
		var ret OrganizationExpanded
		return ret
	}
	return *o.Data
}

// GetDataOk returns a tuple with the Data field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ReadOrganizationExpandedResponse) GetDataOk() (*OrganizationExpanded, bool) {
	if o == nil || IsNil(o.Data) {
		return nil, false
	}
	return o.Data, true
}

// HasData returns a boolean if a field has been set.
func (o *ReadOrganizationExpandedResponse) HasData() bool {
	if o != nil && !IsNil(o.Data) {
		return true
	}

	return false
}

// SetData gets a reference to the given OrganizationExpanded and assigns it to the Data field.
func (o *ReadOrganizationExpandedResponse) SetData(v OrganizationExpanded) {
	o.Data = &v
}

func (o ReadOrganizationExpandedResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ReadOrganizationExpandedResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Data) {
		toSerialize["data"] = o.Data
	}
	return toSerialize, nil
}

type NullableReadOrganizationExpandedResponse struct {
	value *ReadOrganizationExpandedResponse
	isSet bool
}

func (v NullableReadOrganizationExpandedResponse) Get() *ReadOrganizationExpandedResponse {
	return v.value
}

func (v *NullableReadOrganizationExpandedResponse) Set(val *ReadOrganizationExpandedResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableReadOrganizationExpandedResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableReadOrganizationExpandedResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReadOrganizationExpandedResponse(val *ReadOrganizationExpandedResponse) *NullableReadOrganizationExpandedResponse {
	return &NullableReadOrganizationExpandedResponse{value: val, isSet: true}
}

func (v NullableReadOrganizationExpandedResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReadOrganizationExpandedResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


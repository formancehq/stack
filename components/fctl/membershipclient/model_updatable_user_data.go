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

// checks if the UpdatableUserData type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpdatableUserData{}

// UpdatableUserData struct for UpdatableUserData
type UpdatableUserData struct {
	Metadata *map[string]string `json:"metadata,omitempty"`
}

// NewUpdatableUserData instantiates a new UpdatableUserData object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdatableUserData() *UpdatableUserData {
	this := UpdatableUserData{}
	return &this
}

// NewUpdatableUserDataWithDefaults instantiates a new UpdatableUserData object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdatableUserDataWithDefaults() *UpdatableUserData {
	this := UpdatableUserData{}
	return &this
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *UpdatableUserData) GetMetadata() map[string]string {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]string
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdatableUserData) GetMetadataOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Metadata) {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *UpdatableUserData) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]string and assigns it to the Metadata field.
func (o *UpdatableUserData) SetMetadata(v map[string]string) {
	o.Metadata = &v
}

func (o UpdatableUserData) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpdatableUserData) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	return toSerialize, nil
}

type NullableUpdatableUserData struct {
	value *UpdatableUserData
	isSet bool
}

func (v NullableUpdatableUserData) Get() *UpdatableUserData {
	return v.value
}

func (v *NullableUpdatableUserData) Set(val *UpdatableUserData) {
	v.value = val
	v.isSet = true
}

func (v NullableUpdatableUserData) IsSet() bool {
	return v.isSet
}

func (v *NullableUpdatableUserData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpdatableUserData(val *UpdatableUserData) *NullableUpdatableUserData {
	return &NullableUpdatableUserData{value: val, isSet: true}
}

func (v NullableUpdatableUserData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpdatableUserData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



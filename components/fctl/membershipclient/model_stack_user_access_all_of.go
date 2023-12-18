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

// checks if the StackUserAccessAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StackUserAccessAllOf{}

// StackUserAccessAllOf struct for StackUserAccessAllOf
type StackUserAccessAllOf struct {
	Roles []string `json:"roles,omitempty"`
}

// NewStackUserAccessAllOf instantiates a new StackUserAccessAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStackUserAccessAllOf() *StackUserAccessAllOf {
	this := StackUserAccessAllOf{}
	return &this
}

// NewStackUserAccessAllOfWithDefaults instantiates a new StackUserAccessAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStackUserAccessAllOfWithDefaults() *StackUserAccessAllOf {
	this := StackUserAccessAllOf{}
	return &this
}

// GetRoles returns the Roles field value if set, zero value otherwise.
func (o *StackUserAccessAllOf) GetRoles() []string {
	if o == nil || IsNil(o.Roles) {
		var ret []string
		return ret
	}
	return o.Roles
}

// GetRolesOk returns a tuple with the Roles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StackUserAccessAllOf) GetRolesOk() ([]string, bool) {
	if o == nil || IsNil(o.Roles) {
		return nil, false
	}
	return o.Roles, true
}

// HasRoles returns a boolean if a field has been set.
func (o *StackUserAccessAllOf) HasRoles() bool {
	if o != nil && !IsNil(o.Roles) {
		return true
	}

	return false
}

// SetRoles gets a reference to the given []string and assigns it to the Roles field.
func (o *StackUserAccessAllOf) SetRoles(v []string) {
	o.Roles = v
}

func (o StackUserAccessAllOf) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StackUserAccessAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Roles) {
		toSerialize["roles"] = o.Roles
	}
	return toSerialize, nil
}

type NullableStackUserAccessAllOf struct {
	value *StackUserAccessAllOf
	isSet bool
}

func (v NullableStackUserAccessAllOf) Get() *StackUserAccessAllOf {
	return v.value
}

func (v *NullableStackUserAccessAllOf) Set(val *StackUserAccessAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableStackUserAccessAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableStackUserAccessAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStackUserAccessAllOf(val *StackUserAccessAllOf) *NullableStackUserAccessAllOf {
	return &NullableStackUserAccessAllOf{value: val, isSet: true}
}

func (v NullableStackUserAccessAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStackUserAccessAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



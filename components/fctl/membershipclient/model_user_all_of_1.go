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

// checks if the UserAllOf1 type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserAllOf1{}

// UserAllOf1 struct for UserAllOf1
type UserAllOf1 struct {
	Role *SystemRole `json:"role,omitempty"`
}

// NewUserAllOf1 instantiates a new UserAllOf1 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserAllOf1() *UserAllOf1 {
	this := UserAllOf1{}
	return &this
}

// NewUserAllOf1WithDefaults instantiates a new UserAllOf1 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserAllOf1WithDefaults() *UserAllOf1 {
	this := UserAllOf1{}
	return &this
}

// GetRole returns the Role field value if set, zero value otherwise.
func (o *UserAllOf1) GetRole() SystemRole {
	if o == nil || IsNil(o.Role) {
		var ret SystemRole
		return ret
	}
	return *o.Role
}

// GetRoleOk returns a tuple with the Role field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAllOf1) GetRoleOk() (*SystemRole, bool) {
	if o == nil || IsNil(o.Role) {
		return nil, false
	}
	return o.Role, true
}

// HasRole returns a boolean if a field has been set.
func (o *UserAllOf1) HasRole() bool {
	if o != nil && !IsNil(o.Role) {
		return true
	}

	return false
}

// SetRole gets a reference to the given SystemRole and assigns it to the Role field.
func (o *UserAllOf1) SetRole(v SystemRole) {
	o.Role = &v
}

func (o UserAllOf1) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserAllOf1) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Role) {
		toSerialize["role"] = o.Role
	}
	return toSerialize, nil
}

type NullableUserAllOf1 struct {
	value *UserAllOf1
	isSet bool
}

func (v NullableUserAllOf1) Get() *UserAllOf1 {
	return v.value
}

func (v *NullableUserAllOf1) Set(val *UserAllOf1) {
	v.value = val
	v.isSet = true
}

func (v NullableUserAllOf1) IsSet() bool {
	return v.isSet
}

func (v *NullableUserAllOf1) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserAllOf1(val *UserAllOf1) *NullableUserAllOf1 {
	return &NullableUserAllOf1{value: val, isSet: true}
}

func (v NullableUserAllOf1) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserAllOf1) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



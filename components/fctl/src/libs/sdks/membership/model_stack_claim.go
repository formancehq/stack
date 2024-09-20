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

// checks if the StackClaim type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StackClaim{}

// StackClaim struct for StackClaim
type StackClaim struct {
	Id string `json:"id"`
	Role Role `json:"role"`
}

// NewStackClaim instantiates a new StackClaim object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStackClaim(id string, role Role) *StackClaim {
	this := StackClaim{}
	this.Id = id
	this.Role = role
	return &this
}

// NewStackClaimWithDefaults instantiates a new StackClaim object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStackClaimWithDefaults() *StackClaim {
	this := StackClaim{}
	return &this
}

// GetId returns the Id field value
func (o *StackClaim) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *StackClaim) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *StackClaim) SetId(v string) {
	o.Id = v
}

// GetRole returns the Role field value
func (o *StackClaim) GetRole() Role {
	if o == nil {
		var ret Role
		return ret
	}

	return o.Role
}

// GetRoleOk returns a tuple with the Role field value
// and a boolean to check if the value has been set.
func (o *StackClaim) GetRoleOk() (*Role, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Role, true
}

// SetRole sets field value
func (o *StackClaim) SetRole(v Role) {
	o.Role = v
}

func (o StackClaim) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StackClaim) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["role"] = o.Role
	return toSerialize, nil
}

type NullableStackClaim struct {
	value *StackClaim
	isSet bool
}

func (v NullableStackClaim) Get() *StackClaim {
	return v.value
}

func (v *NullableStackClaim) Set(val *StackClaim) {
	v.value = val
	v.isSet = true
}

func (v NullableStackClaim) IsSet() bool {
	return v.isSet
}

func (v *NullableStackClaim) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStackClaim(val *StackClaim) *NullableStackClaim {
	return &NullableStackClaim{value: val, isSet: true}
}

func (v NullableStackClaim) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStackClaim) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



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

// checks if the OrganizationUser type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &OrganizationUser{}

// OrganizationUser struct for OrganizationUser
type OrganizationUser struct {
	Role Role `json:"role"`
	Email string `json:"email"`
	Id string `json:"id"`
}

// NewOrganizationUser instantiates a new OrganizationUser object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOrganizationUser(role Role, email string, id string) *OrganizationUser {
	this := OrganizationUser{}
	this.Role = role
	this.Email = email
	this.Id = id
	return &this
}

// NewOrganizationUserWithDefaults instantiates a new OrganizationUser object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOrganizationUserWithDefaults() *OrganizationUser {
	this := OrganizationUser{}
	return &this
}

// GetRole returns the Role field value
func (o *OrganizationUser) GetRole() Role {
	if o == nil {
		var ret Role
		return ret
	}

	return o.Role
}

// GetRoleOk returns a tuple with the Role field value
// and a boolean to check if the value has been set.
func (o *OrganizationUser) GetRoleOk() (*Role, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Role, true
}

// SetRole sets field value
func (o *OrganizationUser) SetRole(v Role) {
	o.Role = v
}

// GetEmail returns the Email field value
func (o *OrganizationUser) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *OrganizationUser) GetEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *OrganizationUser) SetEmail(v string) {
	o.Email = v
}

// GetId returns the Id field value
func (o *OrganizationUser) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *OrganizationUser) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *OrganizationUser) SetId(v string) {
	o.Id = v
}

func (o OrganizationUser) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o OrganizationUser) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["role"] = o.Role
	toSerialize["email"] = o.Email
	toSerialize["id"] = o.Id
	return toSerialize, nil
}

type NullableOrganizationUser struct {
	value *OrganizationUser
	isSet bool
}

func (v NullableOrganizationUser) Get() *OrganizationUser {
	return v.value
}

func (v *NullableOrganizationUser) Set(val *OrganizationUser) {
	v.value = val
	v.isSet = true
}

func (v NullableOrganizationUser) IsSet() bool {
	return v.isSet
}

func (v *NullableOrganizationUser) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOrganizationUser(val *OrganizationUser) *NullableOrganizationUser {
	return &NullableOrganizationUser{value: val, isSet: true}
}

func (v NullableOrganizationUser) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOrganizationUser) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



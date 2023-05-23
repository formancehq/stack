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

// checks if the UserDataAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserDataAllOf{}

// UserDataAllOf struct for UserDataAllOf
type UserDataAllOf struct {
	Email string `json:"email"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewUserDataAllOf instantiates a new UserDataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserDataAllOf(email string) *UserDataAllOf {
	this := UserDataAllOf{}
	this.Email = email
	return &this
}

// NewUserDataAllOfWithDefaults instantiates a new UserDataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserDataAllOfWithDefaults() *UserDataAllOf {
	this := UserDataAllOf{}
	return &this
}

// GetEmail returns the Email field value
func (o *UserDataAllOf) GetEmail() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Email
}

// GetEmailOk returns a tuple with the Email field value
// and a boolean to check if the value has been set.
func (o *UserDataAllOf) GetEmailOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Email, true
}

// SetEmail sets field value
func (o *UserDataAllOf) SetEmail(v string) {
	o.Email = v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *UserDataAllOf) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserDataAllOf) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *UserDataAllOf) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *UserDataAllOf) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

func (o UserDataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserDataAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["email"] = o.Email
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	return toSerialize, nil
}

type NullableUserDataAllOf struct {
	value *UserDataAllOf
	isSet bool
}

func (v NullableUserDataAllOf) Get() *UserDataAllOf {
	return v.value
}

func (v *NullableUserDataAllOf) Set(val *UserDataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableUserDataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableUserDataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserDataAllOf(val *UserDataAllOf) *NullableUserDataAllOf {
	return &NullableUserDataAllOf{value: val, isSet: true}
}

func (v NullableUserDataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserDataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

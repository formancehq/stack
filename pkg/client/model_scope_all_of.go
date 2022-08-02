/*
Auth API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: AUTH_VERSION
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package authclient

import (
	"encoding/json"
)

// ScopeAllOf struct for ScopeAllOf
type ScopeAllOf struct {
	Id string `json:"id"`
	Transient []string `json:"transient,omitempty"`
}

// NewScopeAllOf instantiates a new ScopeAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewScopeAllOf(id string) *ScopeAllOf {
	this := ScopeAllOf{}
	this.Id = id
	return &this
}

// NewScopeAllOfWithDefaults instantiates a new ScopeAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewScopeAllOfWithDefaults() *ScopeAllOf {
	this := ScopeAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *ScopeAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *ScopeAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *ScopeAllOf) SetId(v string) {
	o.Id = v
}

// GetTransient returns the Transient field value if set, zero value otherwise.
func (o *ScopeAllOf) GetTransient() []string {
	if o == nil || o.Transient == nil {
		var ret []string
		return ret
	}
	return o.Transient
}

// GetTransientOk returns a tuple with the Transient field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScopeAllOf) GetTransientOk() ([]string, bool) {
	if o == nil || o.Transient == nil {
		return nil, false
	}
	return o.Transient, true
}

// HasTransient returns a boolean if a field has been set.
func (o *ScopeAllOf) HasTransient() bool {
	if o != nil && o.Transient != nil {
		return true
	}

	return false
}

// SetTransient gets a reference to the given []string and assigns it to the Transient field.
func (o *ScopeAllOf) SetTransient(v []string) {
	o.Transient = v
}

func (o ScopeAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["id"] = o.Id
	}
	if o.Transient != nil {
		toSerialize["transient"] = o.Transient
	}
	return json.Marshal(toSerialize)
}

type NullableScopeAllOf struct {
	value *ScopeAllOf
	isSet bool
}

func (v NullableScopeAllOf) Get() *ScopeAllOf {
	return v.value
}

func (v *NullableScopeAllOf) Set(val *ScopeAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableScopeAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableScopeAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableScopeAllOf(val *ScopeAllOf) *NullableScopeAllOf {
	return &NullableScopeAllOf{value: val, isSet: true}
}

func (v NullableScopeAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableScopeAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



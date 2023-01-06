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

// TaskDescriptorModulrDescriptor struct for TaskDescriptorModulrDescriptor
type TaskDescriptorModulrDescriptor struct {
	Name *string `json:"name,omitempty"`
	Key *string `json:"key,omitempty"`
	AccountID *string `json:"accountID,omitempty"`
}

// NewTaskDescriptorModulrDescriptor instantiates a new TaskDescriptorModulrDescriptor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskDescriptorModulrDescriptor() *TaskDescriptorModulrDescriptor {
	this := TaskDescriptorModulrDescriptor{}
	return &this
}

// NewTaskDescriptorModulrDescriptorWithDefaults instantiates a new TaskDescriptorModulrDescriptor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskDescriptorModulrDescriptorWithDefaults() *TaskDescriptorModulrDescriptor {
	this := TaskDescriptorModulrDescriptor{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TaskDescriptorModulrDescriptor) GetName() string {
	if o == nil || isNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorModulrDescriptor) GetNameOk() (*string, bool) {
	if o == nil || isNil(o.Name) {
    return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TaskDescriptorModulrDescriptor) HasName() bool {
	if o != nil && !isNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TaskDescriptorModulrDescriptor) SetName(v string) {
	o.Name = &v
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *TaskDescriptorModulrDescriptor) GetKey() string {
	if o == nil || isNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorModulrDescriptor) GetKeyOk() (*string, bool) {
	if o == nil || isNil(o.Key) {
    return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *TaskDescriptorModulrDescriptor) HasKey() bool {
	if o != nil && !isNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *TaskDescriptorModulrDescriptor) SetKey(v string) {
	o.Key = &v
}

// GetAccountID returns the AccountID field value if set, zero value otherwise.
func (o *TaskDescriptorModulrDescriptor) GetAccountID() string {
	if o == nil || isNil(o.AccountID) {
		var ret string
		return ret
	}
	return *o.AccountID
}

// GetAccountIDOk returns a tuple with the AccountID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorModulrDescriptor) GetAccountIDOk() (*string, bool) {
	if o == nil || isNil(o.AccountID) {
    return nil, false
	}
	return o.AccountID, true
}

// HasAccountID returns a boolean if a field has been set.
func (o *TaskDescriptorModulrDescriptor) HasAccountID() bool {
	if o != nil && !isNil(o.AccountID) {
		return true
	}

	return false
}

// SetAccountID gets a reference to the given string and assigns it to the AccountID field.
func (o *TaskDescriptorModulrDescriptor) SetAccountID(v string) {
	o.AccountID = &v
}

func (o TaskDescriptorModulrDescriptor) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !isNil(o.Key) {
		toSerialize["key"] = o.Key
	}
	if !isNil(o.AccountID) {
		toSerialize["accountID"] = o.AccountID
	}
	return json.Marshal(toSerialize)
}

type NullableTaskDescriptorModulrDescriptor struct {
	value *TaskDescriptorModulrDescriptor
	isSet bool
}

func (v NullableTaskDescriptorModulrDescriptor) Get() *TaskDescriptorModulrDescriptor {
	return v.value
}

func (v *NullableTaskDescriptorModulrDescriptor) Set(val *TaskDescriptorModulrDescriptor) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskDescriptorModulrDescriptor) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskDescriptorModulrDescriptor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskDescriptorModulrDescriptor(val *TaskDescriptorModulrDescriptor) *NullableTaskDescriptorModulrDescriptor {
	return &NullableTaskDescriptorModulrDescriptor{value: val, isSet: true}
}

func (v NullableTaskDescriptorModulrDescriptor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskDescriptorModulrDescriptor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

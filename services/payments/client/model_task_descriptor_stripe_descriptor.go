/*
Payments API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// checks if the TaskDescriptorStripeDescriptor type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TaskDescriptorStripeDescriptor{}

// TaskDescriptorStripeDescriptor struct for TaskDescriptorStripeDescriptor
type TaskDescriptorStripeDescriptor struct {
	Name    *string `json:"name,omitempty"`
	Main    *bool   `json:"main,omitempty"`
	Account *string `json:"account,omitempty"`
}

// NewTaskDescriptorStripeDescriptor instantiates a new TaskDescriptorStripeDescriptor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskDescriptorStripeDescriptor() *TaskDescriptorStripeDescriptor {
	this := TaskDescriptorStripeDescriptor{}
	return &this
}

// NewTaskDescriptorStripeDescriptorWithDefaults instantiates a new TaskDescriptorStripeDescriptor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskDescriptorStripeDescriptorWithDefaults() *TaskDescriptorStripeDescriptor {
	this := TaskDescriptorStripeDescriptor{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TaskDescriptorStripeDescriptor) GetName() string {
	if o == nil || isNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorStripeDescriptor) GetNameOk() (*string, bool) {
	if o == nil || isNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TaskDescriptorStripeDescriptor) HasName() bool {
	if o != nil && !isNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TaskDescriptorStripeDescriptor) SetName(v string) {
	o.Name = &v
}

// GetMain returns the Main field value if set, zero value otherwise.
func (o *TaskDescriptorStripeDescriptor) GetMain() bool {
	if o == nil || isNil(o.Main) {
		var ret bool
		return ret
	}
	return *o.Main
}

// GetMainOk returns a tuple with the Main field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorStripeDescriptor) GetMainOk() (*bool, bool) {
	if o == nil || isNil(o.Main) {
		return nil, false
	}
	return o.Main, true
}

// HasMain returns a boolean if a field has been set.
func (o *TaskDescriptorStripeDescriptor) HasMain() bool {
	if o != nil && !isNil(o.Main) {
		return true
	}

	return false
}

// SetMain gets a reference to the given bool and assigns it to the Main field.
func (o *TaskDescriptorStripeDescriptor) SetMain(v bool) {
	o.Main = &v
}

// GetAccount returns the Account field value if set, zero value otherwise.
func (o *TaskDescriptorStripeDescriptor) GetAccount() string {
	if o == nil || isNil(o.Account) {
		var ret string
		return ret
	}
	return *o.Account
}

// GetAccountOk returns a tuple with the Account field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorStripeDescriptor) GetAccountOk() (*string, bool) {
	if o == nil || isNil(o.Account) {
		return nil, false
	}
	return o.Account, true
}

// HasAccount returns a boolean if a field has been set.
func (o *TaskDescriptorStripeDescriptor) HasAccount() bool {
	if o != nil && !isNil(o.Account) {
		return true
	}

	return false
}

// SetAccount gets a reference to the given string and assigns it to the Account field.
func (o *TaskDescriptorStripeDescriptor) SetAccount(v string) {
	o.Account = &v
}

func (o TaskDescriptorStripeDescriptor) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TaskDescriptorStripeDescriptor) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !isNil(o.Main) {
		toSerialize["main"] = o.Main
	}
	if !isNil(o.Account) {
		toSerialize["account"] = o.Account
	}
	return toSerialize, nil
}

type NullableTaskDescriptorStripeDescriptor struct {
	value *TaskDescriptorStripeDescriptor
	isSet bool
}

func (v NullableTaskDescriptorStripeDescriptor) Get() *TaskDescriptorStripeDescriptor {
	return v.value
}

func (v *NullableTaskDescriptorStripeDescriptor) Set(val *TaskDescriptorStripeDescriptor) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskDescriptorStripeDescriptor) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskDescriptorStripeDescriptor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskDescriptorStripeDescriptor(val *TaskDescriptorStripeDescriptor) *NullableTaskDescriptorStripeDescriptor {
	return &NullableTaskDescriptorStripeDescriptor{value: val, isSet: true}
}

func (v NullableTaskDescriptorStripeDescriptor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskDescriptorStripeDescriptor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

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

// checks if the TaskDescriptorWiseDescriptor type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TaskDescriptorWiseDescriptor{}

// TaskDescriptorWiseDescriptor struct for TaskDescriptorWiseDescriptor
type TaskDescriptorWiseDescriptor struct {
	Name      *string `json:"name,omitempty"`
	Key       *string `json:"key,omitempty"`
	ProfileID *int32  `json:"profileID,omitempty"`
}

// NewTaskDescriptorWiseDescriptor instantiates a new TaskDescriptorWiseDescriptor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskDescriptorWiseDescriptor() *TaskDescriptorWiseDescriptor {
	this := TaskDescriptorWiseDescriptor{}
	return &this
}

// NewTaskDescriptorWiseDescriptorWithDefaults instantiates a new TaskDescriptorWiseDescriptor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskDescriptorWiseDescriptorWithDefaults() *TaskDescriptorWiseDescriptor {
	this := TaskDescriptorWiseDescriptor{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TaskDescriptorWiseDescriptor) GetName() string {
	if o == nil || isNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorWiseDescriptor) GetNameOk() (*string, bool) {
	if o == nil || isNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TaskDescriptorWiseDescriptor) HasName() bool {
	if o != nil && !isNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TaskDescriptorWiseDescriptor) SetName(v string) {
	o.Name = &v
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *TaskDescriptorWiseDescriptor) GetKey() string {
	if o == nil || isNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorWiseDescriptor) GetKeyOk() (*string, bool) {
	if o == nil || isNil(o.Key) {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *TaskDescriptorWiseDescriptor) HasKey() bool {
	if o != nil && !isNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *TaskDescriptorWiseDescriptor) SetKey(v string) {
	o.Key = &v
}

// GetProfileID returns the ProfileID field value if set, zero value otherwise.
func (o *TaskDescriptorWiseDescriptor) GetProfileID() int32 {
	if o == nil || isNil(o.ProfileID) {
		var ret int32
		return ret
	}
	return *o.ProfileID
}

// GetProfileIDOk returns a tuple with the ProfileID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorWiseDescriptor) GetProfileIDOk() (*int32, bool) {
	if o == nil || isNil(o.ProfileID) {
		return nil, false
	}
	return o.ProfileID, true
}

// HasProfileID returns a boolean if a field has been set.
func (o *TaskDescriptorWiseDescriptor) HasProfileID() bool {
	if o != nil && !isNil(o.ProfileID) {
		return true
	}

	return false
}

// SetProfileID gets a reference to the given int32 and assigns it to the ProfileID field.
func (o *TaskDescriptorWiseDescriptor) SetProfileID(v int32) {
	o.ProfileID = &v
}

func (o TaskDescriptorWiseDescriptor) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TaskDescriptorWiseDescriptor) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !isNil(o.Key) {
		toSerialize["key"] = o.Key
	}
	if !isNil(o.ProfileID) {
		toSerialize["profileID"] = o.ProfileID
	}
	return toSerialize, nil
}

type NullableTaskDescriptorWiseDescriptor struct {
	value *TaskDescriptorWiseDescriptor
	isSet bool
}

func (v NullableTaskDescriptorWiseDescriptor) Get() *TaskDescriptorWiseDescriptor {
	return v.value
}

func (v *NullableTaskDescriptorWiseDescriptor) Set(val *TaskDescriptorWiseDescriptor) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskDescriptorWiseDescriptor) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskDescriptorWiseDescriptor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskDescriptorWiseDescriptor(val *TaskDescriptorWiseDescriptor) *NullableTaskDescriptorWiseDescriptor {
	return &NullableTaskDescriptorWiseDescriptor{value: val, isSet: true}
}

func (v NullableTaskDescriptorWiseDescriptor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskDescriptorWiseDescriptor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

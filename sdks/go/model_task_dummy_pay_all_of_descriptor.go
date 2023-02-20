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

// checks if the TaskDummyPayAllOfDescriptor type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TaskDummyPayAllOfDescriptor{}

// TaskDummyPayAllOfDescriptor struct for TaskDummyPayAllOfDescriptor
type TaskDummyPayAllOfDescriptor struct {
	Name *string `json:"name,omitempty"`
	Key *string `json:"key,omitempty"`
	FileName *string `json:"fileName,omitempty"`
}

// NewTaskDummyPayAllOfDescriptor instantiates a new TaskDummyPayAllOfDescriptor object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskDummyPayAllOfDescriptor() *TaskDummyPayAllOfDescriptor {
	this := TaskDummyPayAllOfDescriptor{}
	return &this
}

// NewTaskDummyPayAllOfDescriptorWithDefaults instantiates a new TaskDummyPayAllOfDescriptor object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskDummyPayAllOfDescriptorWithDefaults() *TaskDummyPayAllOfDescriptor {
	this := TaskDummyPayAllOfDescriptor{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TaskDummyPayAllOfDescriptor) GetName() string {
	if o == nil || isNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDummyPayAllOfDescriptor) GetNameOk() (*string, bool) {
	if o == nil || isNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TaskDummyPayAllOfDescriptor) HasName() bool {
	if o != nil && !isNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TaskDummyPayAllOfDescriptor) SetName(v string) {
	o.Name = &v
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *TaskDummyPayAllOfDescriptor) GetKey() string {
	if o == nil || isNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDummyPayAllOfDescriptor) GetKeyOk() (*string, bool) {
	if o == nil || isNil(o.Key) {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *TaskDummyPayAllOfDescriptor) HasKey() bool {
	if o != nil && !isNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *TaskDummyPayAllOfDescriptor) SetKey(v string) {
	o.Key = &v
}

// GetFileName returns the FileName field value if set, zero value otherwise.
func (o *TaskDummyPayAllOfDescriptor) GetFileName() string {
	if o == nil || isNil(o.FileName) {
		var ret string
		return ret
	}
	return *o.FileName
}

// GetFileNameOk returns a tuple with the FileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDummyPayAllOfDescriptor) GetFileNameOk() (*string, bool) {
	if o == nil || isNil(o.FileName) {
		return nil, false
	}
	return o.FileName, true
}

// HasFileName returns a boolean if a field has been set.
func (o *TaskDummyPayAllOfDescriptor) HasFileName() bool {
	if o != nil && !isNil(o.FileName) {
		return true
	}

	return false
}

// SetFileName gets a reference to the given string and assigns it to the FileName field.
func (o *TaskDummyPayAllOfDescriptor) SetFileName(v string) {
	o.FileName = &v
}

func (o TaskDummyPayAllOfDescriptor) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TaskDummyPayAllOfDescriptor) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !isNil(o.Key) {
		toSerialize["key"] = o.Key
	}
	if !isNil(o.FileName) {
		toSerialize["fileName"] = o.FileName
	}
	return toSerialize, nil
}

type NullableTaskDummyPayAllOfDescriptor struct {
	value *TaskDummyPayAllOfDescriptor
	isSet bool
}

func (v NullableTaskDummyPayAllOfDescriptor) Get() *TaskDummyPayAllOfDescriptor {
	return v.value
}

func (v *NullableTaskDummyPayAllOfDescriptor) Set(val *TaskDummyPayAllOfDescriptor) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskDummyPayAllOfDescriptor) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskDummyPayAllOfDescriptor) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskDummyPayAllOfDescriptor(val *TaskDummyPayAllOfDescriptor) *NullableTaskDummyPayAllOfDescriptor {
	return &NullableTaskDummyPayAllOfDescriptor{value: val, isSet: true}
}

func (v NullableTaskDummyPayAllOfDescriptor) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskDummyPayAllOfDescriptor) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



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

// TaskResponse struct for TaskResponse
type TaskResponse struct {
	Data TasksCursorCursorAllOfDataInner `json:"data"`
}

// NewTaskResponse instantiates a new TaskResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskResponse(data TasksCursorCursorAllOfDataInner) *TaskResponse {
	this := TaskResponse{}
	this.Data = data
	return &this
}

// NewTaskResponseWithDefaults instantiates a new TaskResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskResponseWithDefaults() *TaskResponse {
	this := TaskResponse{}
	return &this
}

// GetData returns the Data field value
func (o *TaskResponse) GetData() TasksCursorCursorAllOfDataInner {
	if o == nil {
		var ret TasksCursorCursorAllOfDataInner
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *TaskResponse) GetDataOk() (*TasksCursorCursorAllOfDataInner, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *TaskResponse) SetData(v TasksCursorCursorAllOfDataInner) {
	o.Data = v
}

func (o TaskResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableTaskResponse struct {
	value *TaskResponse
	isSet bool
}

func (v NullableTaskResponse) Get() *TaskResponse {
	return v.value
}

func (v *NullableTaskResponse) Set(val *TaskResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskResponse(val *TaskResponse) *NullableTaskResponse {
	return &NullableTaskResponse{value: val, isSet: true}
}

func (v NullableTaskResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



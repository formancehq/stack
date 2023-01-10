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

// ListRunsResponseCursorAllOf struct for ListRunsResponseCursorAllOf
type ListRunsResponseCursorAllOf struct {
	Data []WorkflowOccurrence `json:"data"`
}

// NewListRunsResponseCursorAllOf instantiates a new ListRunsResponseCursorAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListRunsResponseCursorAllOf(data []WorkflowOccurrence) *ListRunsResponseCursorAllOf {
	this := ListRunsResponseCursorAllOf{}
	this.Data = data
	return &this
}

// NewListRunsResponseCursorAllOfWithDefaults instantiates a new ListRunsResponseCursorAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListRunsResponseCursorAllOfWithDefaults() *ListRunsResponseCursorAllOf {
	this := ListRunsResponseCursorAllOf{}
	return &this
}

// GetData returns the Data field value
func (o *ListRunsResponseCursorAllOf) GetData() []WorkflowOccurrence {
	if o == nil {
		var ret []WorkflowOccurrence
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *ListRunsResponseCursorAllOf) GetDataOk() ([]WorkflowOccurrence, bool) {
	if o == nil {
    return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *ListRunsResponseCursorAllOf) SetData(v []WorkflowOccurrence) {
	o.Data = v
}

func (o ListRunsResponseCursorAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableListRunsResponseCursorAllOf struct {
	value *ListRunsResponseCursorAllOf
	isSet bool
}

func (v NullableListRunsResponseCursorAllOf) Get() *ListRunsResponseCursorAllOf {
	return v.value
}

func (v *NullableListRunsResponseCursorAllOf) Set(val *ListRunsResponseCursorAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableListRunsResponseCursorAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableListRunsResponseCursorAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListRunsResponseCursorAllOf(val *ListRunsResponseCursorAllOf) *NullableListRunsResponseCursorAllOf {
	return &NullableListRunsResponseCursorAllOf{value: val, isSet: true}
}

func (v NullableListRunsResponseCursorAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListRunsResponseCursorAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



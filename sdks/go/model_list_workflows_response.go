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

// ListWorkflowsResponse struct for ListWorkflowsResponse
type ListWorkflowsResponse struct {
	Cursor ListWorkflowsResponseCursor `json:"cursor"`
}

// NewListWorkflowsResponse instantiates a new ListWorkflowsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListWorkflowsResponse(cursor ListWorkflowsResponseCursor) *ListWorkflowsResponse {
	this := ListWorkflowsResponse{}
	this.Cursor = cursor
	return &this
}

// NewListWorkflowsResponseWithDefaults instantiates a new ListWorkflowsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListWorkflowsResponseWithDefaults() *ListWorkflowsResponse {
	this := ListWorkflowsResponse{}
	return &this
}

// GetCursor returns the Cursor field value
func (o *ListWorkflowsResponse) GetCursor() ListWorkflowsResponseCursor {
	if o == nil {
		var ret ListWorkflowsResponseCursor
		return ret
	}

	return o.Cursor
}

// GetCursorOk returns a tuple with the Cursor field value
// and a boolean to check if the value has been set.
func (o *ListWorkflowsResponse) GetCursorOk() (*ListWorkflowsResponseCursor, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Cursor, true
}

// SetCursor sets field value
func (o *ListWorkflowsResponse) SetCursor(v ListWorkflowsResponseCursor) {
	o.Cursor = v
}

func (o ListWorkflowsResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["cursor"] = o.Cursor
	}
	return json.Marshal(toSerialize)
}

type NullableListWorkflowsResponse struct {
	value *ListWorkflowsResponse
	isSet bool
}

func (v NullableListWorkflowsResponse) Get() *ListWorkflowsResponse {
	return v.value
}

func (v *NullableListWorkflowsResponse) Set(val *ListWorkflowsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableListWorkflowsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableListWorkflowsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListWorkflowsResponse(val *ListWorkflowsResponse) *NullableListWorkflowsResponse {
	return &NullableListWorkflowsResponse{value: val, isSet: true}
}

func (v NullableListWorkflowsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListWorkflowsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



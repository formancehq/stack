/*
Formance Stack API

Open, modular foundation for unique payments flows  # Introduction This API is documented in **OpenAPI format**.  # Authentication Formance Stack offers one forms of authentication:   - OAuth2 OAuth2 - an open protocol to allow secure authorization in a simple and standard method from web, mobile and desktop applications. <SecurityDefinitions /> 

API version: v1.0.20230301
Contact: support@formance.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package formance

import (
	"encoding/json"
)

// checks if the ListRunsResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ListRunsResponse{}

// ListRunsResponse struct for ListRunsResponse
type ListRunsResponse struct {
	Data []WorkflowInstance `json:"data"`
}

// NewListRunsResponse instantiates a new ListRunsResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewListRunsResponse(data []WorkflowInstance) *ListRunsResponse {
	this := ListRunsResponse{}
	this.Data = data
	return &this
}

// NewListRunsResponseWithDefaults instantiates a new ListRunsResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewListRunsResponseWithDefaults() *ListRunsResponse {
	this := ListRunsResponse{}
	return &this
}

// GetData returns the Data field value
func (o *ListRunsResponse) GetData() []WorkflowInstance {
	if o == nil {
		var ret []WorkflowInstance
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *ListRunsResponse) GetDataOk() ([]WorkflowInstance, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *ListRunsResponse) SetData(v []WorkflowInstance) {
	o.Data = v
}

func (o ListRunsResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ListRunsResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableListRunsResponse struct {
	value *ListRunsResponse
	isSet bool
}

func (v NullableListRunsResponse) Get() *ListRunsResponse {
	return v.value
}

func (v *NullableListRunsResponse) Set(val *ListRunsResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableListRunsResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableListRunsResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableListRunsResponse(val *ListRunsResponse) *NullableListRunsResponse {
	return &NullableListRunsResponse{value: val, isSet: true}
}

func (v NullableListRunsResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableListRunsResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



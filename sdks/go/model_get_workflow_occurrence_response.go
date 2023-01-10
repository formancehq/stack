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

// GetWorkflowOccurrenceResponse struct for GetWorkflowOccurrenceResponse
type GetWorkflowOccurrenceResponse struct {
	Data WorkflowOccurrence `json:"data"`
}

// NewGetWorkflowOccurrenceResponse instantiates a new GetWorkflowOccurrenceResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetWorkflowOccurrenceResponse(data WorkflowOccurrence) *GetWorkflowOccurrenceResponse {
	this := GetWorkflowOccurrenceResponse{}
	this.Data = data
	return &this
}

// NewGetWorkflowOccurrenceResponseWithDefaults instantiates a new GetWorkflowOccurrenceResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetWorkflowOccurrenceResponseWithDefaults() *GetWorkflowOccurrenceResponse {
	this := GetWorkflowOccurrenceResponse{}
	return &this
}

// GetData returns the Data field value
func (o *GetWorkflowOccurrenceResponse) GetData() WorkflowOccurrence {
	if o == nil {
		var ret WorkflowOccurrence
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *GetWorkflowOccurrenceResponse) GetDataOk() (*WorkflowOccurrence, bool) {
	if o == nil {
    return nil, false
	}
	return &o.Data, true
}

// SetData sets field value
func (o *GetWorkflowOccurrenceResponse) SetData(v WorkflowOccurrence) {
	o.Data = v
}

func (o GetWorkflowOccurrenceResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["data"] = o.Data
	}
	return json.Marshal(toSerialize)
}

type NullableGetWorkflowOccurrenceResponse struct {
	value *GetWorkflowOccurrenceResponse
	isSet bool
}

func (v NullableGetWorkflowOccurrenceResponse) Get() *GetWorkflowOccurrenceResponse {
	return v.value
}

func (v *NullableGetWorkflowOccurrenceResponse) Set(val *GetWorkflowOccurrenceResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetWorkflowOccurrenceResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetWorkflowOccurrenceResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetWorkflowOccurrenceResponse(val *GetWorkflowOccurrenceResponse) *NullableGetWorkflowOccurrenceResponse {
	return &NullableGetWorkflowOccurrenceResponse{value: val, isSet: true}
}

func (v NullableGetWorkflowOccurrenceResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetWorkflowOccurrenceResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



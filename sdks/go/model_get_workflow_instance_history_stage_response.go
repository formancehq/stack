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

// checks if the GetWorkflowInstanceHistoryStageResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetWorkflowInstanceHistoryStageResponse{}

// GetWorkflowInstanceHistoryStageResponse struct for GetWorkflowInstanceHistoryStageResponse
type GetWorkflowInstanceHistoryStageResponse struct {
	Data []WorkflowInstanceHistoryStage `json:"data"`
}

// NewGetWorkflowInstanceHistoryStageResponse instantiates a new GetWorkflowInstanceHistoryStageResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetWorkflowInstanceHistoryStageResponse(data []WorkflowInstanceHistoryStage) *GetWorkflowInstanceHistoryStageResponse {
	this := GetWorkflowInstanceHistoryStageResponse{}
	this.Data = data
	return &this
}

// NewGetWorkflowInstanceHistoryStageResponseWithDefaults instantiates a new GetWorkflowInstanceHistoryStageResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetWorkflowInstanceHistoryStageResponseWithDefaults() *GetWorkflowInstanceHistoryStageResponse {
	this := GetWorkflowInstanceHistoryStageResponse{}
	return &this
}

// GetData returns the Data field value
func (o *GetWorkflowInstanceHistoryStageResponse) GetData() []WorkflowInstanceHistoryStage {
	if o == nil {
		var ret []WorkflowInstanceHistoryStage
		return ret
	}

	return o.Data
}

// GetDataOk returns a tuple with the Data field value
// and a boolean to check if the value has been set.
func (o *GetWorkflowInstanceHistoryStageResponse) GetDataOk() ([]WorkflowInstanceHistoryStage, bool) {
	if o == nil {
		return nil, false
	}
	return o.Data, true
}

// SetData sets field value
func (o *GetWorkflowInstanceHistoryStageResponse) SetData(v []WorkflowInstanceHistoryStage) {
	o.Data = v
}

func (o GetWorkflowInstanceHistoryStageResponse) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetWorkflowInstanceHistoryStageResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["data"] = o.Data
	return toSerialize, nil
}

type NullableGetWorkflowInstanceHistoryStageResponse struct {
	value *GetWorkflowInstanceHistoryStageResponse
	isSet bool
}

func (v NullableGetWorkflowInstanceHistoryStageResponse) Get() *GetWorkflowInstanceHistoryStageResponse {
	return v.value
}

func (v *NullableGetWorkflowInstanceHistoryStageResponse) Set(val *GetWorkflowInstanceHistoryStageResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableGetWorkflowInstanceHistoryStageResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableGetWorkflowInstanceHistoryStageResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetWorkflowInstanceHistoryStageResponse(val *GetWorkflowInstanceHistoryStageResponse) *NullableGetWorkflowInstanceHistoryStageResponse {
	return &NullableGetWorkflowInstanceHistoryStageResponse{value: val, isSet: true}
}

func (v NullableGetWorkflowInstanceHistoryStageResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetWorkflowInstanceHistoryStageResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



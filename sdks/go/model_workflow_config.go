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

// checks if the WorkflowConfig type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WorkflowConfig{}

// WorkflowConfig struct for WorkflowConfig
type WorkflowConfig struct {
	Name *string `json:"name,omitempty"`
	Stages []map[string]interface{} `json:"stages"`
}

// NewWorkflowConfig instantiates a new WorkflowConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkflowConfig(stages []map[string]interface{}) *WorkflowConfig {
	this := WorkflowConfig{}
	this.Stages = stages
	return &this
}

// NewWorkflowConfigWithDefaults instantiates a new WorkflowConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkflowConfigWithDefaults() *WorkflowConfig {
	this := WorkflowConfig{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *WorkflowConfig) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowConfig) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *WorkflowConfig) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *WorkflowConfig) SetName(v string) {
	o.Name = &v
}

// GetStages returns the Stages field value
func (o *WorkflowConfig) GetStages() []map[string]interface{} {
	if o == nil {
		var ret []map[string]interface{}
		return ret
	}

	return o.Stages
}

// GetStagesOk returns a tuple with the Stages field value
// and a boolean to check if the value has been set.
func (o *WorkflowConfig) GetStagesOk() ([]map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}
	return o.Stages, true
}

// SetStages sets field value
func (o *WorkflowConfig) SetStages(v []map[string]interface{}) {
	o.Stages = v
}

func (o WorkflowConfig) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WorkflowConfig) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	toSerialize["stages"] = o.Stages
	return toSerialize, nil
}

type NullableWorkflowConfig struct {
	value *WorkflowConfig
	isSet bool
}

func (v NullableWorkflowConfig) Get() *WorkflowConfig {
	return v.value
}

func (v *NullableWorkflowConfig) Set(val *WorkflowConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkflowConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkflowConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkflowConfig(val *WorkflowConfig) *NullableWorkflowConfig {
	return &NullableWorkflowConfig{value: val, isSet: true}
}

func (v NullableWorkflowConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkflowConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



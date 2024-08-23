/*
Formance Simple ingester Service API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ingesterclient

import (
	"encoding/json"
	"time"
)

// checks if the Pipeline type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Pipeline{}

// Pipeline struct for Pipeline
type Pipeline struct {
	Module string `json:"module"`
	ConnectorID string `json:"connectorID"`
	Id string `json:"id"`
	State State `json:"state"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewPipeline instantiates a new Pipeline object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipeline(module string, connectorID string, id string, state State, createdAt time.Time) *Pipeline {
	this := Pipeline{}
	this.Module = module
	this.ConnectorID = connectorID
	this.Id = id
	this.State = state
	this.CreatedAt = createdAt
	return &this
}

// NewPipelineWithDefaults instantiates a new Pipeline object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelineWithDefaults() *Pipeline {
	this := Pipeline{}
	return &this
}

// GetModule returns the Module field value
func (o *Pipeline) GetModule() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Module
}

// GetModuleOk returns a tuple with the Module field value
// and a boolean to check if the value has been set.
func (o *Pipeline) GetModuleOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Module, true
}

// SetModule sets field value
func (o *Pipeline) SetModule(v string) {
	o.Module = v
}

// GetConnectorID returns the ConnectorID field value
func (o *Pipeline) GetConnectorID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ConnectorID
}

// GetConnectorIDOk returns a tuple with the ConnectorID field value
// and a boolean to check if the value has been set.
func (o *Pipeline) GetConnectorIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ConnectorID, true
}

// SetConnectorID sets field value
func (o *Pipeline) SetConnectorID(v string) {
	o.ConnectorID = v
}

// GetId returns the Id field value
func (o *Pipeline) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Pipeline) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Pipeline) SetId(v string) {
	o.Id = v
}

// GetState returns the State field value
func (o *Pipeline) GetState() State {
	if o == nil {
		var ret State
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *Pipeline) GetStateOk() (*State, bool) {
	if o == nil {
		return nil, false
	}
	return &o.State, true
}

// SetState sets field value
func (o *Pipeline) SetState(v State) {
	o.State = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *Pipeline) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *Pipeline) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *Pipeline) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

func (o Pipeline) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Pipeline) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["module"] = o.Module
	toSerialize["connectorID"] = o.ConnectorID
	toSerialize["id"] = o.Id
	toSerialize["state"] = o.State
	toSerialize["createdAt"] = o.CreatedAt
	return toSerialize, nil
}

type NullablePipeline struct {
	value *Pipeline
	isSet bool
}

func (v NullablePipeline) Get() *Pipeline {
	return v.value
}

func (v *NullablePipeline) Set(val *Pipeline) {
	v.value = val
	v.isSet = true
}

func (v NullablePipeline) IsSet() bool {
	return v.isSet
}

func (v *NullablePipeline) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipeline(val *Pipeline) *NullablePipeline {
	return &NullablePipeline{value: val, isSet: true}
}

func (v NullablePipeline) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipeline) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



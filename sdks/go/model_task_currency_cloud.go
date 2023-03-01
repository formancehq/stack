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
	"time"
)

// checks if the TaskCurrencyCloud type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TaskCurrencyCloud{}

// TaskCurrencyCloud struct for TaskCurrencyCloud
type TaskCurrencyCloud struct {
	Id string `json:"id"`
	ConnectorId string `json:"connectorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Status PaymentStatus `json:"status"`
	State map[string]interface{} `json:"state"`
	Error *string `json:"error,omitempty"`
	Descriptor TaskCurrencyCloudAllOfDescriptor `json:"descriptor"`
}

// NewTaskCurrencyCloud instantiates a new TaskCurrencyCloud object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskCurrencyCloud(id string, connectorId string, createdAt time.Time, updatedAt time.Time, status PaymentStatus, state map[string]interface{}, descriptor TaskCurrencyCloudAllOfDescriptor) *TaskCurrencyCloud {
	this := TaskCurrencyCloud{}
	this.Id = id
	this.ConnectorId = connectorId
	this.CreatedAt = createdAt
	this.UpdatedAt = updatedAt
	this.Status = status
	this.State = state
	this.Descriptor = descriptor
	return &this
}

// NewTaskCurrencyCloudWithDefaults instantiates a new TaskCurrencyCloud object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskCurrencyCloudWithDefaults() *TaskCurrencyCloud {
	this := TaskCurrencyCloud{}
	return &this
}

// GetId returns the Id field value
func (o *TaskCurrencyCloud) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *TaskCurrencyCloud) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *TaskCurrencyCloud) SetId(v string) {
	o.Id = v
}

// GetConnectorId returns the ConnectorId field value
func (o *TaskCurrencyCloud) GetConnectorId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ConnectorId
}

// GetConnectorIdOk returns a tuple with the ConnectorId field value
// and a boolean to check if the value has been set.
func (o *TaskCurrencyCloud) GetConnectorIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ConnectorId, true
}

// SetConnectorId sets field value
func (o *TaskCurrencyCloud) SetConnectorId(v string) {
	o.ConnectorId = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *TaskCurrencyCloud) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *TaskCurrencyCloud) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *TaskCurrencyCloud) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

// GetUpdatedAt returns the UpdatedAt field value
func (o *TaskCurrencyCloud) GetUpdatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value
// and a boolean to check if the value has been set.
func (o *TaskCurrencyCloud) GetUpdatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UpdatedAt, true
}

// SetUpdatedAt sets field value
func (o *TaskCurrencyCloud) SetUpdatedAt(v time.Time) {
	o.UpdatedAt = v
}

// GetStatus returns the Status field value
func (o *TaskCurrencyCloud) GetStatus() PaymentStatus {
	if o == nil {
		var ret PaymentStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *TaskCurrencyCloud) GetStatusOk() (*PaymentStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *TaskCurrencyCloud) SetStatus(v PaymentStatus) {
	o.Status = v
}

// GetState returns the State field value
func (o *TaskCurrencyCloud) GetState() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *TaskCurrencyCloud) GetStateOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.State, true
}

// SetState sets field value
func (o *TaskCurrencyCloud) SetState(v map[string]interface{}) {
	o.State = v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *TaskCurrencyCloud) GetError() string {
	if o == nil || IsNil(o.Error) {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskCurrencyCloud) GetErrorOk() (*string, bool) {
	if o == nil || IsNil(o.Error) {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *TaskCurrencyCloud) HasError() bool {
	if o != nil && !IsNil(o.Error) {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *TaskCurrencyCloud) SetError(v string) {
	o.Error = &v
}

// GetDescriptor returns the Descriptor field value
func (o *TaskCurrencyCloud) GetDescriptor() TaskCurrencyCloudAllOfDescriptor {
	if o == nil {
		var ret TaskCurrencyCloudAllOfDescriptor
		return ret
	}

	return o.Descriptor
}

// GetDescriptorOk returns a tuple with the Descriptor field value
// and a boolean to check if the value has been set.
func (o *TaskCurrencyCloud) GetDescriptorOk() (*TaskCurrencyCloudAllOfDescriptor, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Descriptor, true
}

// SetDescriptor sets field value
func (o *TaskCurrencyCloud) SetDescriptor(v TaskCurrencyCloudAllOfDescriptor) {
	o.Descriptor = v
}

func (o TaskCurrencyCloud) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TaskCurrencyCloud) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["connectorId"] = o.ConnectorId
	toSerialize["createdAt"] = o.CreatedAt
	toSerialize["updatedAt"] = o.UpdatedAt
	toSerialize["status"] = o.Status
	toSerialize["state"] = o.State
	if !IsNil(o.Error) {
		toSerialize["error"] = o.Error
	}
	toSerialize["descriptor"] = o.Descriptor
	return toSerialize, nil
}

type NullableTaskCurrencyCloud struct {
	value *TaskCurrencyCloud
	isSet bool
}

func (v NullableTaskCurrencyCloud) Get() *TaskCurrencyCloud {
	return v.value
}

func (v *NullableTaskCurrencyCloud) Set(val *TaskCurrencyCloud) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskCurrencyCloud) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskCurrencyCloud) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskCurrencyCloud(val *TaskCurrencyCloud) *NullableTaskCurrencyCloud {
	return &NullableTaskCurrencyCloud{value: val, isSet: true}
}

func (v NullableTaskCurrencyCloud) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskCurrencyCloud) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



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
	"time"
)

// checks if the TaskStripe type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TaskStripe{}

// TaskStripe struct for TaskStripe
type TaskStripe struct {
	Id string `json:"id"`
	ConnectorId string `json:"connectorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Status PaymentStatus `json:"status"`
	State map[string]interface{} `json:"state"`
	Error *string `json:"error,omitempty"`
	Descriptor TaskStripeAllOfDescriptor `json:"descriptor"`
}

// NewTaskStripe instantiates a new TaskStripe object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskStripe(id string, connectorId string, createdAt time.Time, updatedAt time.Time, status PaymentStatus, state map[string]interface{}, descriptor TaskStripeAllOfDescriptor) *TaskStripe {
	this := TaskStripe{}
	this.Id = id
	this.ConnectorId = connectorId
	this.CreatedAt = createdAt
	this.UpdatedAt = updatedAt
	this.Status = status
	this.State = state
	this.Descriptor = descriptor
	return &this
}

// NewTaskStripeWithDefaults instantiates a new TaskStripe object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskStripeWithDefaults() *TaskStripe {
	this := TaskStripe{}
	return &this
}

// GetId returns the Id field value
func (o *TaskStripe) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *TaskStripe) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *TaskStripe) SetId(v string) {
	o.Id = v
}

// GetConnectorId returns the ConnectorId field value
func (o *TaskStripe) GetConnectorId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ConnectorId
}

// GetConnectorIdOk returns a tuple with the ConnectorId field value
// and a boolean to check if the value has been set.
func (o *TaskStripe) GetConnectorIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ConnectorId, true
}

// SetConnectorId sets field value
func (o *TaskStripe) SetConnectorId(v string) {
	o.ConnectorId = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *TaskStripe) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *TaskStripe) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *TaskStripe) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

// GetUpdatedAt returns the UpdatedAt field value
func (o *TaskStripe) GetUpdatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value
// and a boolean to check if the value has been set.
func (o *TaskStripe) GetUpdatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UpdatedAt, true
}

// SetUpdatedAt sets field value
func (o *TaskStripe) SetUpdatedAt(v time.Time) {
	o.UpdatedAt = v
}

// GetStatus returns the Status field value
func (o *TaskStripe) GetStatus() PaymentStatus {
	if o == nil {
		var ret PaymentStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *TaskStripe) GetStatusOk() (*PaymentStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *TaskStripe) SetStatus(v PaymentStatus) {
	o.Status = v
}

// GetState returns the State field value
func (o *TaskStripe) GetState() map[string]interface{} {
	if o == nil {
		var ret map[string]interface{}
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *TaskStripe) GetStateOk() (map[string]interface{}, bool) {
	if o == nil {
		return map[string]interface{}{}, false
	}
	return o.State, true
}

// SetState sets field value
func (o *TaskStripe) SetState(v map[string]interface{}) {
	o.State = v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *TaskStripe) GetError() string {
	if o == nil || isNil(o.Error) {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskStripe) GetErrorOk() (*string, bool) {
	if o == nil || isNil(o.Error) {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *TaskStripe) HasError() bool {
	if o != nil && !isNil(o.Error) {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *TaskStripe) SetError(v string) {
	o.Error = &v
}

// GetDescriptor returns the Descriptor field value
func (o *TaskStripe) GetDescriptor() TaskStripeAllOfDescriptor {
	if o == nil {
		var ret TaskStripeAllOfDescriptor
		return ret
	}

	return o.Descriptor
}

// GetDescriptorOk returns a tuple with the Descriptor field value
// and a boolean to check if the value has been set.
func (o *TaskStripe) GetDescriptorOk() (*TaskStripeAllOfDescriptor, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Descriptor, true
}

// SetDescriptor sets field value
func (o *TaskStripe) SetDescriptor(v TaskStripeAllOfDescriptor) {
	o.Descriptor = v
}

func (o TaskStripe) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TaskStripe) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["connectorId"] = o.ConnectorId
	toSerialize["createdAt"] = o.CreatedAt
	toSerialize["updatedAt"] = o.UpdatedAt
	toSerialize["status"] = o.Status
	toSerialize["state"] = o.State
	if !isNil(o.Error) {
		toSerialize["error"] = o.Error
	}
	toSerialize["descriptor"] = o.Descriptor
	return toSerialize, nil
}

type NullableTaskStripe struct {
	value *TaskStripe
	isSet bool
}

func (v NullableTaskStripe) Get() *TaskStripe {
	return v.value
}

func (v *NullableTaskStripe) Set(val *TaskStripe) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskStripe) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskStripe) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskStripe(val *TaskStripe) *NullableTaskStripe {
	return &NullableTaskStripe{value: val, isSet: true}
}

func (v NullableTaskStripe) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskStripe) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



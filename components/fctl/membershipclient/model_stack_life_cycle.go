/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package membershipclient

import (
	"encoding/json"
	"time"
)

// checks if the StackLifeCycle type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StackLifeCycle{}

// StackLifeCycle struct for StackLifeCycle
type StackLifeCycle struct {
	Status string `json:"status"`
	State string `json:"state"`
	ExpectedStatus string `json:"expectedStatus"`
	LastStateUpdate time.Time `json:"lastStateUpdate"`
	LastExpectedStatusUpdate time.Time `json:"lastExpectedStatusUpdate"`
	LastStatusUpdate time.Time `json:"lastStatusUpdate"`
}

// NewStackLifeCycle instantiates a new StackLifeCycle object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStackLifeCycle(status string, state string, expectedStatus string, lastStateUpdate time.Time, lastExpectedStatusUpdate time.Time, lastStatusUpdate time.Time) *StackLifeCycle {
	this := StackLifeCycle{}
	this.Status = status
	this.State = state
	this.ExpectedStatus = expectedStatus
	this.LastStateUpdate = lastStateUpdate
	this.LastExpectedStatusUpdate = lastExpectedStatusUpdate
	this.LastStatusUpdate = lastStatusUpdate
	return &this
}

// NewStackLifeCycleWithDefaults instantiates a new StackLifeCycle object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStackLifeCycleWithDefaults() *StackLifeCycle {
	this := StackLifeCycle{}
	return &this
}

// GetStatus returns the Status field value
func (o *StackLifeCycle) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *StackLifeCycle) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *StackLifeCycle) SetStatus(v string) {
	o.Status = v
}

// GetState returns the State field value
func (o *StackLifeCycle) GetState() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.State
}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
func (o *StackLifeCycle) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.State, true
}

// SetState sets field value
func (o *StackLifeCycle) SetState(v string) {
	o.State = v
}

// GetExpectedStatus returns the ExpectedStatus field value
func (o *StackLifeCycle) GetExpectedStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ExpectedStatus
}

// GetExpectedStatusOk returns a tuple with the ExpectedStatus field value
// and a boolean to check if the value has been set.
func (o *StackLifeCycle) GetExpectedStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ExpectedStatus, true
}

// SetExpectedStatus sets field value
func (o *StackLifeCycle) SetExpectedStatus(v string) {
	o.ExpectedStatus = v
}

// GetLastStateUpdate returns the LastStateUpdate field value
func (o *StackLifeCycle) GetLastStateUpdate() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.LastStateUpdate
}

// GetLastStateUpdateOk returns a tuple with the LastStateUpdate field value
// and a boolean to check if the value has been set.
func (o *StackLifeCycle) GetLastStateUpdateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LastStateUpdate, true
}

// SetLastStateUpdate sets field value
func (o *StackLifeCycle) SetLastStateUpdate(v time.Time) {
	o.LastStateUpdate = v
}

// GetLastExpectedStatusUpdate returns the LastExpectedStatusUpdate field value
func (o *StackLifeCycle) GetLastExpectedStatusUpdate() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.LastExpectedStatusUpdate
}

// GetLastExpectedStatusUpdateOk returns a tuple with the LastExpectedStatusUpdate field value
// and a boolean to check if the value has been set.
func (o *StackLifeCycle) GetLastExpectedStatusUpdateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LastExpectedStatusUpdate, true
}

// SetLastExpectedStatusUpdate sets field value
func (o *StackLifeCycle) SetLastExpectedStatusUpdate(v time.Time) {
	o.LastExpectedStatusUpdate = v
}

// GetLastStatusUpdate returns the LastStatusUpdate field value
func (o *StackLifeCycle) GetLastStatusUpdate() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.LastStatusUpdate
}

// GetLastStatusUpdateOk returns a tuple with the LastStatusUpdate field value
// and a boolean to check if the value has been set.
func (o *StackLifeCycle) GetLastStatusUpdateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LastStatusUpdate, true
}

// SetLastStatusUpdate sets field value
func (o *StackLifeCycle) SetLastStatusUpdate(v time.Time) {
	o.LastStatusUpdate = v
}

func (o StackLifeCycle) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StackLifeCycle) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["status"] = o.Status
	toSerialize["state"] = o.State
	toSerialize["expectedStatus"] = o.ExpectedStatus
	toSerialize["lastStateUpdate"] = o.LastStateUpdate
	toSerialize["lastExpectedStatusUpdate"] = o.LastExpectedStatusUpdate
	toSerialize["lastStatusUpdate"] = o.LastStatusUpdate
	return toSerialize, nil
}

type NullableStackLifeCycle struct {
	value *StackLifeCycle
	isSet bool
}

func (v NullableStackLifeCycle) Get() *StackLifeCycle {
	return v.value
}

func (v *NullableStackLifeCycle) Set(val *StackLifeCycle) {
	v.value = val
	v.isSet = true
}

func (v NullableStackLifeCycle) IsSet() bool {
	return v.isSet
}

func (v *NullableStackLifeCycle) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStackLifeCycle(val *StackLifeCycle) *NullableStackLifeCycle {
	return &NullableStackLifeCycle{value: val, isSet: true}
}

func (v NullableStackLifeCycle) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStackLifeCycle) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



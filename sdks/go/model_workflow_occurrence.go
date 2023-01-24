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

// checks if the WorkflowOccurrence type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WorkflowOccurrence{}

// WorkflowOccurrence struct for WorkflowOccurrence
type WorkflowOccurrence struct {
	WorkflowID string `json:"workflowID"`
	Id string `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Statuses []StageStatus `json:"statuses"`
}

// NewWorkflowOccurrence instantiates a new WorkflowOccurrence object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkflowOccurrence(workflowID string, id string, createdAt time.Time, updatedAt time.Time, statuses []StageStatus) *WorkflowOccurrence {
	this := WorkflowOccurrence{}
	this.WorkflowID = workflowID
	this.Id = id
	this.CreatedAt = createdAt
	this.UpdatedAt = updatedAt
	this.Statuses = statuses
	return &this
}

// NewWorkflowOccurrenceWithDefaults instantiates a new WorkflowOccurrence object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkflowOccurrenceWithDefaults() *WorkflowOccurrence {
	this := WorkflowOccurrence{}
	return &this
}

// GetWorkflowID returns the WorkflowID field value
func (o *WorkflowOccurrence) GetWorkflowID() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.WorkflowID
}

// GetWorkflowIDOk returns a tuple with the WorkflowID field value
// and a boolean to check if the value has been set.
func (o *WorkflowOccurrence) GetWorkflowIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.WorkflowID, true
}

// SetWorkflowID sets field value
func (o *WorkflowOccurrence) SetWorkflowID(v string) {
	o.WorkflowID = v
}

// GetId returns the Id field value
func (o *WorkflowOccurrence) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WorkflowOccurrence) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WorkflowOccurrence) SetId(v string) {
	o.Id = v
}

// GetCreatedAt returns the CreatedAt field value
func (o *WorkflowOccurrence) GetCreatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *WorkflowOccurrence) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *WorkflowOccurrence) SetCreatedAt(v time.Time) {
	o.CreatedAt = v
}

// GetUpdatedAt returns the UpdatedAt field value
func (o *WorkflowOccurrence) GetUpdatedAt() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value
// and a boolean to check if the value has been set.
func (o *WorkflowOccurrence) GetUpdatedAtOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UpdatedAt, true
}

// SetUpdatedAt sets field value
func (o *WorkflowOccurrence) SetUpdatedAt(v time.Time) {
	o.UpdatedAt = v
}

// GetStatuses returns the Statuses field value
func (o *WorkflowOccurrence) GetStatuses() []StageStatus {
	if o == nil {
		var ret []StageStatus
		return ret
	}

	return o.Statuses
}

// GetStatusesOk returns a tuple with the Statuses field value
// and a boolean to check if the value has been set.
func (o *WorkflowOccurrence) GetStatusesOk() ([]StageStatus, bool) {
	if o == nil {
		return nil, false
	}
	return o.Statuses, true
}

// SetStatuses sets field value
func (o *WorkflowOccurrence) SetStatuses(v []StageStatus) {
	o.Statuses = v
}

func (o WorkflowOccurrence) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WorkflowOccurrence) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["workflowID"] = o.WorkflowID
	toSerialize["id"] = o.Id
	toSerialize["createdAt"] = o.CreatedAt
	toSerialize["updatedAt"] = o.UpdatedAt
	toSerialize["statuses"] = o.Statuses
	return toSerialize, nil
}

type NullableWorkflowOccurrence struct {
	value *WorkflowOccurrence
	isSet bool
}

func (v NullableWorkflowOccurrence) Get() *WorkflowOccurrence {
	return v.value
}

func (v *NullableWorkflowOccurrence) Set(val *WorkflowOccurrence) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkflowOccurrence) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkflowOccurrence) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkflowOccurrence(val *WorkflowOccurrence) *NullableWorkflowOccurrence {
	return &NullableWorkflowOccurrence{value: val, isSet: true}
}

func (v NullableWorkflowOccurrence) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkflowOccurrence) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



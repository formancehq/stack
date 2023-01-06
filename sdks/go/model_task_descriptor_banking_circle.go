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

// TaskDescriptorBankingCircle struct for TaskDescriptorBankingCircle
type TaskDescriptorBankingCircle struct {
	// The connector code
	Provider *string `json:"provider,omitempty"`
	// The date when the task was created
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// The task status
	Status *string `json:"status,omitempty"`
	// The error message if the task failed
	Error *string `json:"error,omitempty"`
	// The task state
	State      map[string]interface{}                 `json:"state,omitempty"`
	Descriptor *TaskDescriptorBankingCircleDescriptor `json:"descriptor,omitempty"`
}

// NewTaskDescriptorBankingCircle instantiates a new TaskDescriptorBankingCircle object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTaskDescriptorBankingCircle() *TaskDescriptorBankingCircle {
	this := TaskDescriptorBankingCircle{}
	return &this
}

// NewTaskDescriptorBankingCircleWithDefaults instantiates a new TaskDescriptorBankingCircle object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTaskDescriptorBankingCircleWithDefaults() *TaskDescriptorBankingCircle {
	this := TaskDescriptorBankingCircle{}
	return &this
}

// GetProvider returns the Provider field value if set, zero value otherwise.
func (o *TaskDescriptorBankingCircle) GetProvider() string {
	if o == nil || isNil(o.Provider) {
		var ret string
		return ret
	}
	return *o.Provider
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorBankingCircle) GetProviderOk() (*string, bool) {
	if o == nil || isNil(o.Provider) {
		return nil, false
	}
	return o.Provider, true
}

// HasProvider returns a boolean if a field has been set.
func (o *TaskDescriptorBankingCircle) HasProvider() bool {
	if o != nil && !isNil(o.Provider) {
		return true
	}

	return false
}

// SetProvider gets a reference to the given string and assigns it to the Provider field.
func (o *TaskDescriptorBankingCircle) SetProvider(v string) {
	o.Provider = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *TaskDescriptorBankingCircle) GetCreatedAt() time.Time {
	if o == nil || isNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorBankingCircle) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || isNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *TaskDescriptorBankingCircle) HasCreatedAt() bool {
	if o != nil && !isNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *TaskDescriptorBankingCircle) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *TaskDescriptorBankingCircle) GetStatus() string {
	if o == nil || isNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorBankingCircle) GetStatusOk() (*string, bool) {
	if o == nil || isNil(o.Status) {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *TaskDescriptorBankingCircle) HasStatus() bool {
	if o != nil && !isNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *TaskDescriptorBankingCircle) SetStatus(v string) {
	o.Status = &v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *TaskDescriptorBankingCircle) GetError() string {
	if o == nil || isNil(o.Error) {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorBankingCircle) GetErrorOk() (*string, bool) {
	if o == nil || isNil(o.Error) {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *TaskDescriptorBankingCircle) HasError() bool {
	if o != nil && !isNil(o.Error) {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *TaskDescriptorBankingCircle) SetError(v string) {
	o.Error = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *TaskDescriptorBankingCircle) GetState() map[string]interface{} {
	if o == nil || isNil(o.State) {
		var ret map[string]interface{}
		return ret
	}
	return o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorBankingCircle) GetStateOk() (map[string]interface{}, bool) {
	if o == nil || isNil(o.State) {
		return map[string]interface{}{}, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *TaskDescriptorBankingCircle) HasState() bool {
	if o != nil && !isNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given map[string]interface{} and assigns it to the State field.
func (o *TaskDescriptorBankingCircle) SetState(v map[string]interface{}) {
	o.State = v
}

// GetDescriptor returns the Descriptor field value if set, zero value otherwise.
func (o *TaskDescriptorBankingCircle) GetDescriptor() TaskDescriptorBankingCircleDescriptor {
	if o == nil || isNil(o.Descriptor) {
		var ret TaskDescriptorBankingCircleDescriptor
		return ret
	}
	return *o.Descriptor
}

// GetDescriptorOk returns a tuple with the Descriptor field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TaskDescriptorBankingCircle) GetDescriptorOk() (*TaskDescriptorBankingCircleDescriptor, bool) {
	if o == nil || isNil(o.Descriptor) {
		return nil, false
	}
	return o.Descriptor, true
}

// HasDescriptor returns a boolean if a field has been set.
func (o *TaskDescriptorBankingCircle) HasDescriptor() bool {
	if o != nil && !isNil(o.Descriptor) {
		return true
	}

	return false
}

// SetDescriptor gets a reference to the given TaskDescriptorBankingCircleDescriptor and assigns it to the Descriptor field.
func (o *TaskDescriptorBankingCircle) SetDescriptor(v TaskDescriptorBankingCircleDescriptor) {
	o.Descriptor = &v
}

func (o TaskDescriptorBankingCircle) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Provider) {
		toSerialize["provider"] = o.Provider
	}
	if !isNil(o.CreatedAt) {
		toSerialize["createdAt"] = o.CreatedAt
	}
	if !isNil(o.Status) {
		toSerialize["status"] = o.Status
	}
	if !isNil(o.Error) {
		toSerialize["error"] = o.Error
	}
	if !isNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !isNil(o.Descriptor) {
		toSerialize["descriptor"] = o.Descriptor
	}
	return json.Marshal(toSerialize)
}

type NullableTaskDescriptorBankingCircle struct {
	value *TaskDescriptorBankingCircle
	isSet bool
}

func (v NullableTaskDescriptorBankingCircle) Get() *TaskDescriptorBankingCircle {
	return v.value
}

func (v *NullableTaskDescriptorBankingCircle) Set(val *TaskDescriptorBankingCircle) {
	v.value = val
	v.isSet = true
}

func (v NullableTaskDescriptorBankingCircle) IsSet() bool {
	return v.isSet
}

func (v *NullableTaskDescriptorBankingCircle) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTaskDescriptorBankingCircle(val *TaskDescriptorBankingCircle) *NullableTaskDescriptorBankingCircle {
	return &NullableTaskDescriptorBankingCircle{value: val, isSet: true}
}

func (v NullableTaskDescriptorBankingCircle) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTaskDescriptorBankingCircle) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

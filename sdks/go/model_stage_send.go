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

// checks if the StageSend type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &StageSend{}

// StageSend struct for StageSend
type StageSend struct {
	Amount *Monetary `json:"amount,omitempty"`
	Destination *StageSendDestination `json:"destination,omitempty"`
	Source *StageSendSource `json:"source,omitempty"`
}

// NewStageSend instantiates a new StageSend object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStageSend() *StageSend {
	this := StageSend{}
	return &this
}

// NewStageSendWithDefaults instantiates a new StageSend object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStageSendWithDefaults() *StageSend {
	this := StageSend{}
	return &this
}

// GetAmount returns the Amount field value if set, zero value otherwise.
func (o *StageSend) GetAmount() Monetary {
	if o == nil || isNil(o.Amount) {
		var ret Monetary
		return ret
	}
	return *o.Amount
}

// GetAmountOk returns a tuple with the Amount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSend) GetAmountOk() (*Monetary, bool) {
	if o == nil || isNil(o.Amount) {
		return nil, false
	}
	return o.Amount, true
}

// HasAmount returns a boolean if a field has been set.
func (o *StageSend) HasAmount() bool {
	if o != nil && !isNil(o.Amount) {
		return true
	}

	return false
}

// SetAmount gets a reference to the given Monetary and assigns it to the Amount field.
func (o *StageSend) SetAmount(v Monetary) {
	o.Amount = &v
}

// GetDestination returns the Destination field value if set, zero value otherwise.
func (o *StageSend) GetDestination() StageSendDestination {
	if o == nil || isNil(o.Destination) {
		var ret StageSendDestination
		return ret
	}
	return *o.Destination
}

// GetDestinationOk returns a tuple with the Destination field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSend) GetDestinationOk() (*StageSendDestination, bool) {
	if o == nil || isNil(o.Destination) {
		return nil, false
	}
	return o.Destination, true
}

// HasDestination returns a boolean if a field has been set.
func (o *StageSend) HasDestination() bool {
	if o != nil && !isNil(o.Destination) {
		return true
	}

	return false
}

// SetDestination gets a reference to the given StageSendDestination and assigns it to the Destination field.
func (o *StageSend) SetDestination(v StageSendDestination) {
	o.Destination = &v
}

// GetSource returns the Source field value if set, zero value otherwise.
func (o *StageSend) GetSource() StageSendSource {
	if o == nil || isNil(o.Source) {
		var ret StageSendSource
		return ret
	}
	return *o.Source
}

// GetSourceOk returns a tuple with the Source field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StageSend) GetSourceOk() (*StageSendSource, bool) {
	if o == nil || isNil(o.Source) {
		return nil, false
	}
	return o.Source, true
}

// HasSource returns a boolean if a field has been set.
func (o *StageSend) HasSource() bool {
	if o != nil && !isNil(o.Source) {
		return true
	}

	return false
}

// SetSource gets a reference to the given StageSendSource and assigns it to the Source field.
func (o *StageSend) SetSource(v StageSendSource) {
	o.Source = &v
}

func (o StageSend) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o StageSend) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Amount) {
		toSerialize["amount"] = o.Amount
	}
	if !isNil(o.Destination) {
		toSerialize["destination"] = o.Destination
	}
	if !isNil(o.Source) {
		toSerialize["source"] = o.Source
	}
	return toSerialize, nil
}

type NullableStageSend struct {
	value *StageSend
	isSet bool
}

func (v NullableStageSend) Get() *StageSend {
	return v.value
}

func (v *NullableStageSend) Set(val *StageSend) {
	v.value = val
	v.isSet = true
}

func (v NullableStageSend) IsSet() bool {
	return v.isSet
}

func (v *NullableStageSend) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStageSend(val *StageSend) *NullableStageSend {
	return &NullableStageSend{value: val, isSet: true}
}

func (v NullableStageSend) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStageSend) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



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

// PaymentMetadataChangelog struct for PaymentMetadataChangelog
type PaymentMetadataChangelog struct {
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Value *string `json:"value,omitempty"`
}

// NewPaymentMetadataChangelog instantiates a new PaymentMetadataChangelog object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaymentMetadataChangelog() *PaymentMetadataChangelog {
	this := PaymentMetadataChangelog{}
	return &this
}

// NewPaymentMetadataChangelogWithDefaults instantiates a new PaymentMetadataChangelog object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaymentMetadataChangelogWithDefaults() *PaymentMetadataChangelog {
	this := PaymentMetadataChangelog{}
	return &this
}

// GetTimestamp returns the Timestamp field value if set, zero value otherwise.
func (o *PaymentMetadataChangelog) GetTimestamp() time.Time {
	if o == nil || isNil(o.Timestamp) {
		var ret time.Time
		return ret
	}
	return *o.Timestamp
}

// GetTimestampOk returns a tuple with the Timestamp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaymentMetadataChangelog) GetTimestampOk() (*time.Time, bool) {
	if o == nil || isNil(o.Timestamp) {
    return nil, false
	}
	return o.Timestamp, true
}

// HasTimestamp returns a boolean if a field has been set.
func (o *PaymentMetadataChangelog) HasTimestamp() bool {
	if o != nil && !isNil(o.Timestamp) {
		return true
	}

	return false
}

// SetTimestamp gets a reference to the given time.Time and assigns it to the Timestamp field.
func (o *PaymentMetadataChangelog) SetTimestamp(v time.Time) {
	o.Timestamp = &v
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *PaymentMetadataChangelog) GetValue() string {
	if o == nil || isNil(o.Value) {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PaymentMetadataChangelog) GetValueOk() (*string, bool) {
	if o == nil || isNil(o.Value) {
    return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *PaymentMetadataChangelog) HasValue() bool {
	if o != nil && !isNil(o.Value) {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *PaymentMetadataChangelog) SetValue(v string) {
	o.Value = &v
}

func (o PaymentMetadataChangelog) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Timestamp) {
		toSerialize["timestamp"] = o.Timestamp
	}
	if !isNil(o.Value) {
		toSerialize["value"] = o.Value
	}
	return json.Marshal(toSerialize)
}

type NullablePaymentMetadataChangelog struct {
	value *PaymentMetadataChangelog
	isSet bool
}

func (v NullablePaymentMetadataChangelog) Get() *PaymentMetadataChangelog {
	return v.value
}

func (v *NullablePaymentMetadataChangelog) Set(val *PaymentMetadataChangelog) {
	v.value = val
	v.isSet = true
}

func (v NullablePaymentMetadataChangelog) IsSet() bool {
	return v.isSet
}

func (v *NullablePaymentMetadataChangelog) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePaymentMetadataChangelog(val *PaymentMetadataChangelog) *NullablePaymentMetadataChangelog {
	return &NullablePaymentMetadataChangelog{value: val, isSet: true}
}

func (v NullablePaymentMetadataChangelog) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePaymentMetadataChangelog) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



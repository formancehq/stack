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

// checks if the PaymentSource type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PaymentSource{}

// PaymentSource struct for PaymentSource
type PaymentSource struct {
	Id string `json:"id"`
}

// NewPaymentSource instantiates a new PaymentSource object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPaymentSource(id string) *PaymentSource {
	this := PaymentSource{}
	this.Id = id
	return &this
}

// NewPaymentSourceWithDefaults instantiates a new PaymentSource object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaymentSourceWithDefaults() *PaymentSource {
	this := PaymentSource{}
	return &this
}

// GetId returns the Id field value
func (o *PaymentSource) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *PaymentSource) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *PaymentSource) SetId(v string) {
	o.Id = v
}

func (o PaymentSource) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PaymentSource) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	return toSerialize, nil
}

type NullablePaymentSource struct {
	value *PaymentSource
	isSet bool
}

func (v NullablePaymentSource) Get() *PaymentSource {
	return v.value
}

func (v *NullablePaymentSource) Set(val *PaymentSource) {
	v.value = val
	v.isSet = true
}

func (v NullablePaymentSource) IsSet() bool {
	return v.isSet
}

func (v *NullablePaymentSource) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePaymentSource(val *PaymentSource) *NullablePaymentSource {
	return &NullablePaymentSource{value: val, isSet: true}
}

func (v NullablePaymentSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePaymentSource) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



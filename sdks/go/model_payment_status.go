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
	"fmt"
)

// PaymentStatus the model 'PaymentStatus'
type PaymentStatus string

// List of PaymentStatus
const (
	PENDING PaymentStatus = "PENDING"
	ACTIVE PaymentStatus = "ACTIVE"
	TERMINATED PaymentStatus = "TERMINATED"
	FAILED PaymentStatus = "FAILED"
)

// All allowed values of PaymentStatus enum
var AllowedPaymentStatusEnumValues = []PaymentStatus{
	"PENDING",
	"ACTIVE",
	"TERMINATED",
	"FAILED",
}

func (v *PaymentStatus) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := PaymentStatus(value)
	for _, existing := range AllowedPaymentStatusEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid PaymentStatus", value)
}

// NewPaymentStatusFromValue returns a pointer to a valid PaymentStatus
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewPaymentStatusFromValue(v string) (*PaymentStatus, error) {
	ev := PaymentStatus(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for PaymentStatus: valid values are %v", v, AllowedPaymentStatusEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v PaymentStatus) IsValid() bool {
	for _, existing := range AllowedPaymentStatusEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to PaymentStatus value
func (v PaymentStatus) Ptr() *PaymentStatus {
	return &v
}

type NullablePaymentStatus struct {
	value *PaymentStatus
	isSet bool
}

func (v NullablePaymentStatus) Get() *PaymentStatus {
	return v.value
}

func (v *NullablePaymentStatus) Set(val *PaymentStatus) {
	v.value = val
	v.isSet = true
}

func (v NullablePaymentStatus) IsSet() bool {
	return v.isSet
}

func (v *NullablePaymentStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePaymentStatus(val *PaymentStatus) *NullablePaymentStatus {
	return &NullablePaymentStatus{value: val, isSet: true}
}

func (v NullablePaymentStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePaymentStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


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

// Connectors the model 'Connectors'
type Connectors string

// List of Connectors
const (
	STRIPE Connectors = "STRIPE"
	DUMMY_PAY Connectors = "DUMMY-PAY"
	SIE Connectors = "SIE"
	MODULR Connectors = "MODULR"
	CURRENCY_CLOUD Connectors = "CURRENCY-CLOUD"
	BANKING_CIRCLE Connectors = "BANKING-CIRCLE"
)

// All allowed values of Connectors enum
var AllowedConnectorsEnumValues = []Connectors{
	"STRIPE",
	"DUMMY-PAY",
	"SIE",
	"MODULR",
	"CURRENCY-CLOUD",
	"BANKING-CIRCLE",
}

func (v *Connectors) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Connectors(value)
	for _, existing := range AllowedConnectorsEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid Connectors", value)
}

// NewConnectorsFromValue returns a pointer to a valid Connectors
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewConnectorsFromValue(v string) (*Connectors, error) {
	ev := Connectors(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for Connectors: valid values are %v", v, AllowedConnectorsEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Connectors) IsValid() bool {
	for _, existing := range AllowedConnectorsEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to Connectors value
func (v Connectors) Ptr() *Connectors {
	return &v
}

type NullableConnectors struct {
	value *Connectors
	isSet bool
}

func (v NullableConnectors) Get() *Connectors {
	return v.value
}

func (v *NullableConnectors) Set(val *Connectors) {
	v.value = val
	v.isSet = true
}

func (v NullableConnectors) IsSet() bool {
	return v.isSet
}

func (v *NullableConnectors) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConnectors(val *Connectors) *NullableConnectors {
	return &NullableConnectors{value: val, isSet: true}
}

func (v NullableConnectors) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConnectors) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


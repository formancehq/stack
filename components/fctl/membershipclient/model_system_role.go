/*
Membership API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package membershipclient

import (
	"encoding/json"
	"fmt"
)

// SystemRole the model 'SystemRole'
type SystemRole string

// List of SystemRole
const (
	USER SystemRole = "USER"
	SYSTEM SystemRole = "SYSTEM"
)

// All allowed values of SystemRole enum
var AllowedSystemRoleEnumValues = []SystemRole{
	"USER",
	"SYSTEM",
}

func (v *SystemRole) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := SystemRole(value)
	for _, existing := range AllowedSystemRoleEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid SystemRole", value)
}

// NewSystemRoleFromValue returns a pointer to a valid SystemRole
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewSystemRoleFromValue(v string) (*SystemRole, error) {
	ev := SystemRole(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for SystemRole: valid values are %v", v, AllowedSystemRoleEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v SystemRole) IsValid() bool {
	for _, existing := range AllowedSystemRoleEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to SystemRole value
func (v SystemRole) Ptr() *SystemRole {
	return &v
}

type NullableSystemRole struct {
	value *SystemRole
	isSet bool
}

func (v NullableSystemRole) Get() *SystemRole {
	return v.value
}

func (v *NullableSystemRole) Set(val *SystemRole) {
	v.value = val
	v.isSet = true
}

func (v NullableSystemRole) IsSet() bool {
	return v.isSet
}

func (v *NullableSystemRole) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSystemRole(val *SystemRole) *NullableSystemRole {
	return &NullableSystemRole{value: val, isSet: true}
}

func (v NullableSystemRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSystemRole) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

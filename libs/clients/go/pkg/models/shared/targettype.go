// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type TargetType string

const (
	TargetTypeTransaction TargetType = "TRANSACTION"
	TargetTypeAccount     TargetType = "ACCOUNT"
)

func (e TargetType) ToPointer() *TargetType {
	return &e
}

func (e *TargetType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "TRANSACTION":
		fallthrough
	case "ACCOUNT":
		*e = TargetType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for TargetType: %v", v)
	}
}

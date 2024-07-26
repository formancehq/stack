// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type V2TargetType string

const (
	V2TargetTypeTransaction V2TargetType = "TRANSACTION"
	V2TargetTypeAccount     V2TargetType = "ACCOUNT"
)

func (e V2TargetType) ToPointer() *V2TargetType {
	return &e
}
func (e *V2TargetType) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "TRANSACTION":
		fallthrough
	case "ACCOUNT":
		*e = V2TargetType(v)
		return nil
	default:
		return fmt.Errorf("invalid value for V2TargetType: %v", v)
	}
}

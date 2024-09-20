// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package sdkerrors

import (
	"encoding/json"
	"fmt"
)

type SchemasErrorCode string

const (
	SchemasErrorCodeValidation SchemasErrorCode = "VALIDATION"
	SchemasErrorCodeNotFound   SchemasErrorCode = "NOT_FOUND"
	SchemasErrorCodeInternal   SchemasErrorCode = "INTERNAL"
)

func (e SchemasErrorCode) ToPointer() *SchemasErrorCode {
	return &e
}
func (e *SchemasErrorCode) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "VALIDATION":
		fallthrough
	case "NOT_FOUND":
		fallthrough
	case "INTERNAL":
		*e = SchemasErrorCode(v)
		return nil
	default:
		return fmt.Errorf("invalid value for SchemasErrorCode: %v", v)
	}
}

// V2Error - General error
type V2Error struct {
	ErrorCode    SchemasErrorCode `json:"errorCode"`
	ErrorMessage string           `json:"errorMessage"`
}

var _ error = &V2Error{}

func (e *V2Error) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

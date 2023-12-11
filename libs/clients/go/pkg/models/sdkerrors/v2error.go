// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package sdkerrors

import (
	"encoding/json"
	"fmt"
)

type V2ErrorErrorCode string

const (
	V2ErrorErrorCodeValidation V2ErrorErrorCode = "VALIDATION"
	V2ErrorErrorCodeNotFound   V2ErrorErrorCode = "NOT_FOUND"
)

func (e V2ErrorErrorCode) ToPointer() *V2ErrorErrorCode {
	return &e
}

func (e *V2ErrorErrorCode) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "VALIDATION":
		fallthrough
	case "NOT_FOUND":
		*e = V2ErrorErrorCode(v)
		return nil
	default:
		return fmt.Errorf("invalid value for V2ErrorErrorCode: %v", v)
	}
}

type V2Error struct {
	ErrorCode    V2ErrorErrorCode `json:"errorCode"`
	ErrorMessage string           `json:"errorMessage"`
}

var _ error = &V2Error{}

func (e *V2Error) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

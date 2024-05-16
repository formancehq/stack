// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"encoding/json"
	"fmt"
)

type Status string

const (
	StatusWaitingForValidation Status = "WAITING_FOR_VALIDATION"
	StatusProcessing           Status = "PROCESSING"
	StatusProcessed            Status = "PROCESSED"
	StatusFailed               Status = "FAILED"
	StatusRejected             Status = "REJECTED"
	StatusValidated            Status = "VALIDATED"
)

func (e Status) ToPointer() *Status {
	return &e
}
func (e *Status) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v {
	case "WAITING_FOR_VALIDATION":
		fallthrough
	case "PROCESSING":
		fallthrough
	case "PROCESSED":
		fallthrough
	case "FAILED":
		fallthrough
	case "REJECTED":
		fallthrough
	case "VALIDATED":
		*e = Status(v)
		return nil
	default:
		return fmt.Errorf("invalid value for Status: %v", v)
	}
}

type UpdateTransferInitiationStatusRequest struct {
	Status Status `json:"status"`
}

func (o *UpdateTransferInitiationStatusRequest) GetStatus() Status {
	if o == nil {
		return Status("")
	}
	return o.Status
}

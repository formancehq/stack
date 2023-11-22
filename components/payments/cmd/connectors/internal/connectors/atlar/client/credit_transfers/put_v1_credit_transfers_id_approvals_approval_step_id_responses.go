// Code generated by go-swagger; DO NOT EDIT.

package credit_transfers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/models"
)

// PutV1CreditTransfersIDApprovalsApprovalStepIDReader is a Reader for the PutV1CreditTransfersIDApprovalsApprovalStepID structure.
type PutV1CreditTransfersIDApprovalsApprovalStepIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPutV1CreditTransfersIDApprovalsApprovalStepIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[PUT /v1/credit-transfers/{id}/approvals/{approvalStepId}] PutV1CreditTransfersIDApprovalsApprovalStepID", response, response.Code())
	}
}

// NewPutV1CreditTransfersIDApprovalsApprovalStepIDOK creates a PutV1CreditTransfersIDApprovalsApprovalStepIDOK with default headers values
func NewPutV1CreditTransfersIDApprovalsApprovalStepIDOK() *PutV1CreditTransfersIDApprovalsApprovalStepIDOK {
	return &PutV1CreditTransfersIDApprovalsApprovalStepIDOK{}
}

/*
PutV1CreditTransfersIDApprovalsApprovalStepIDOK describes a response with status code 200, with default header values.

the, now, approved identified CreditTransfer
*/
type PutV1CreditTransfersIDApprovalsApprovalStepIDOK struct {
	Payload *models.Payment
}

// IsSuccess returns true when this put v1 credit transfers Id approvals approval step Id o k response has a 2xx status code
func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this put v1 credit transfers Id approvals approval step Id o k response has a 3xx status code
func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this put v1 credit transfers Id approvals approval step Id o k response has a 4xx status code
func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this put v1 credit transfers Id approvals approval step Id o k response has a 5xx status code
func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this put v1 credit transfers Id approvals approval step Id o k response a status code equal to that given
func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the put v1 credit transfers Id approvals approval step Id o k response
func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) Code() int {
	return 200
}

func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) Error() string {
	return fmt.Sprintf("[PUT /v1/credit-transfers/{id}/approvals/{approvalStepId}][%d] putV1CreditTransfersIdApprovalsApprovalStepIdOK  %+v", 200, o.Payload)
}

func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) String() string {
	return fmt.Sprintf("[PUT /v1/credit-transfers/{id}/approvals/{approvalStepId}][%d] putV1CreditTransfersIdApprovalsApprovalStepIdOK  %+v", 200, o.Payload)
}

func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) GetPayload() *models.Payment {
	return o.Payload
}

func (o *PutV1CreditTransfersIDApprovalsApprovalStepIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Payment)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

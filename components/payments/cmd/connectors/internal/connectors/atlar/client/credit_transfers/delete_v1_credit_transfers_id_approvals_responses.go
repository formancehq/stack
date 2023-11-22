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

// DeleteV1CreditTransfersIDApprovalsReader is a Reader for the DeleteV1CreditTransfersIDApprovals structure.
type DeleteV1CreditTransfersIDApprovalsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteV1CreditTransfersIDApprovalsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteV1CreditTransfersIDApprovalsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[DELETE /v1/credit-transfers/{id}/approvals] DeleteV1CreditTransfersIDApprovals", response, response.Code())
	}
}

// NewDeleteV1CreditTransfersIDApprovalsOK creates a DeleteV1CreditTransfersIDApprovalsOK with default headers values
func NewDeleteV1CreditTransfersIDApprovalsOK() *DeleteV1CreditTransfersIDApprovalsOK {
	return &DeleteV1CreditTransfersIDApprovalsOK{}
}

/*
DeleteV1CreditTransfersIDApprovalsOK describes a response with status code 200, with default header values.

the, now, rejected identified CreditTransfer
*/
type DeleteV1CreditTransfersIDApprovalsOK struct {
	Payload *models.Payment
}

// IsSuccess returns true when this delete v1 credit transfers Id approvals o k response has a 2xx status code
func (o *DeleteV1CreditTransfersIDApprovalsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete v1 credit transfers Id approvals o k response has a 3xx status code
func (o *DeleteV1CreditTransfersIDApprovalsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete v1 credit transfers Id approvals o k response has a 4xx status code
func (o *DeleteV1CreditTransfersIDApprovalsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete v1 credit transfers Id approvals o k response has a 5xx status code
func (o *DeleteV1CreditTransfersIDApprovalsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this delete v1 credit transfers Id approvals o k response a status code equal to that given
func (o *DeleteV1CreditTransfersIDApprovalsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the delete v1 credit transfers Id approvals o k response
func (o *DeleteV1CreditTransfersIDApprovalsOK) Code() int {
	return 200
}

func (o *DeleteV1CreditTransfersIDApprovalsOK) Error() string {
	return fmt.Sprintf("[DELETE /v1/credit-transfers/{id}/approvals][%d] deleteV1CreditTransfersIdApprovalsOK  %+v", 200, o.Payload)
}

func (o *DeleteV1CreditTransfersIDApprovalsOK) String() string {
	return fmt.Sprintf("[DELETE /v1/credit-transfers/{id}/approvals][%d] deleteV1CreditTransfersIdApprovalsOK  %+v", 200, o.Payload)
}

func (o *DeleteV1CreditTransfersIDApprovalsOK) GetPayload() *models.Payment {
	return o.Payload
}

func (o *DeleteV1CreditTransfersIDApprovalsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Payment)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

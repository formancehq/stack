// Code generated by go-swagger; DO NOT EDIT.

package accounts

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/models"
)

// GetV1AccountsIDReader is a Reader for the GetV1AccountsID structure.
type GetV1AccountsIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetV1AccountsIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetV1AccountsIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetV1AccountsIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /v1/accounts/{id}] GetV1AccountsID", response, response.Code())
	}
}

// NewGetV1AccountsIDOK creates a GetV1AccountsIDOK with default headers values
func NewGetV1AccountsIDOK() *GetV1AccountsIDOK {
	return &GetV1AccountsIDOK{}
}

/*
GetV1AccountsIDOK describes a response with status code 200, with default header values.

desc
*/
type GetV1AccountsIDOK struct {
	Payload *models.Account
}

// IsSuccess returns true when this get v1 accounts Id o k response has a 2xx status code
func (o *GetV1AccountsIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get v1 accounts Id o k response has a 3xx status code
func (o *GetV1AccountsIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 accounts Id o k response has a 4xx status code
func (o *GetV1AccountsIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get v1 accounts Id o k response has a 5xx status code
func (o *GetV1AccountsIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get v1 accounts Id o k response a status code equal to that given
func (o *GetV1AccountsIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get v1 accounts Id o k response
func (o *GetV1AccountsIDOK) Code() int {
	return 200
}

func (o *GetV1AccountsIDOK) Error() string {
	return fmt.Sprintf("[GET /v1/accounts/{id}][%d] getV1AccountsIdOK  %+v", 200, o.Payload)
}

func (o *GetV1AccountsIDOK) String() string {
	return fmt.Sprintf("[GET /v1/accounts/{id}][%d] getV1AccountsIdOK  %+v", 200, o.Payload)
}

func (o *GetV1AccountsIDOK) GetPayload() *models.Account {
	return o.Payload
}

func (o *GetV1AccountsIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Account)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetV1AccountsIDNotFound creates a GetV1AccountsIDNotFound with default headers values
func NewGetV1AccountsIDNotFound() *GetV1AccountsIDNotFound {
	return &GetV1AccountsIDNotFound{}
}

/*
GetV1AccountsIDNotFound describes a response with status code 404, with default header values.

identified account doesn't exist
*/
type GetV1AccountsIDNotFound struct {
	Payload string
}

// IsSuccess returns true when this get v1 accounts Id not found response has a 2xx status code
func (o *GetV1AccountsIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get v1 accounts Id not found response has a 3xx status code
func (o *GetV1AccountsIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 accounts Id not found response has a 4xx status code
func (o *GetV1AccountsIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get v1 accounts Id not found response has a 5xx status code
func (o *GetV1AccountsIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get v1 accounts Id not found response a status code equal to that given
func (o *GetV1AccountsIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get v1 accounts Id not found response
func (o *GetV1AccountsIDNotFound) Code() int {
	return 404
}

func (o *GetV1AccountsIDNotFound) Error() string {
	return fmt.Sprintf("[GET /v1/accounts/{id}][%d] getV1AccountsIdNotFound  %+v", 404, o.Payload)
}

func (o *GetV1AccountsIDNotFound) String() string {
	return fmt.Sprintf("[GET /v1/accounts/{id}][%d] getV1AccountsIdNotFound  %+v", 404, o.Payload)
}

func (o *GetV1AccountsIDNotFound) GetPayload() string {
	return o.Payload
}

func (o *GetV1AccountsIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

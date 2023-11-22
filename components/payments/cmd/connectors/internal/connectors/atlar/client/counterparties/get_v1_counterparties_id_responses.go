// Code generated by go-swagger; DO NOT EDIT.

package counterparties

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/models"
)

// GetV1CounterpartiesIDReader is a Reader for the GetV1CounterpartiesID structure.
type GetV1CounterpartiesIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetV1CounterpartiesIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetV1CounterpartiesIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetV1CounterpartiesIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /v1/counterparties/{id}] GetV1CounterpartiesID", response, response.Code())
	}
}

// NewGetV1CounterpartiesIDOK creates a GetV1CounterpartiesIDOK with default headers values
func NewGetV1CounterpartiesIDOK() *GetV1CounterpartiesIDOK {
	return &GetV1CounterpartiesIDOK{}
}

/*
GetV1CounterpartiesIDOK describes a response with status code 200, with default header values.

The identified counterparty.
*/
type GetV1CounterpartiesIDOK struct {
	Payload *models.Counterparty
}

// IsSuccess returns true when this get v1 counterparties Id o k response has a 2xx status code
func (o *GetV1CounterpartiesIDOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get v1 counterparties Id o k response has a 3xx status code
func (o *GetV1CounterpartiesIDOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 counterparties Id o k response has a 4xx status code
func (o *GetV1CounterpartiesIDOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get v1 counterparties Id o k response has a 5xx status code
func (o *GetV1CounterpartiesIDOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get v1 counterparties Id o k response a status code equal to that given
func (o *GetV1CounterpartiesIDOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get v1 counterparties Id o k response
func (o *GetV1CounterpartiesIDOK) Code() int {
	return 200
}

func (o *GetV1CounterpartiesIDOK) Error() string {
	return fmt.Sprintf("[GET /v1/counterparties/{id}][%d] getV1CounterpartiesIdOK  %+v", 200, o.Payload)
}

func (o *GetV1CounterpartiesIDOK) String() string {
	return fmt.Sprintf("[GET /v1/counterparties/{id}][%d] getV1CounterpartiesIdOK  %+v", 200, o.Payload)
}

func (o *GetV1CounterpartiesIDOK) GetPayload() *models.Counterparty {
	return o.Payload
}

func (o *GetV1CounterpartiesIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Counterparty)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetV1CounterpartiesIDNotFound creates a GetV1CounterpartiesIDNotFound with default headers values
func NewGetV1CounterpartiesIDNotFound() *GetV1CounterpartiesIDNotFound {
	return &GetV1CounterpartiesIDNotFound{}
}

/*
GetV1CounterpartiesIDNotFound describes a response with status code 404, with default header values.

The identified counterparty doesn't exist.
*/
type GetV1CounterpartiesIDNotFound struct {
	Payload string
}

// IsSuccess returns true when this get v1 counterparties Id not found response has a 2xx status code
func (o *GetV1CounterpartiesIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get v1 counterparties Id not found response has a 3xx status code
func (o *GetV1CounterpartiesIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 counterparties Id not found response has a 4xx status code
func (o *GetV1CounterpartiesIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get v1 counterparties Id not found response has a 5xx status code
func (o *GetV1CounterpartiesIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get v1 counterparties Id not found response a status code equal to that given
func (o *GetV1CounterpartiesIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get v1 counterparties Id not found response
func (o *GetV1CounterpartiesIDNotFound) Code() int {
	return 404
}

func (o *GetV1CounterpartiesIDNotFound) Error() string {
	return fmt.Sprintf("[GET /v1/counterparties/{id}][%d] getV1CounterpartiesIdNotFound  %+v", 404, o.Payload)
}

func (o *GetV1CounterpartiesIDNotFound) String() string {
	return fmt.Sprintf("[GET /v1/counterparties/{id}][%d] getV1CounterpartiesIdNotFound  %+v", 404, o.Payload)
}

func (o *GetV1CounterpartiesIDNotFound) GetPayload() string {
	return o.Payload
}

func (o *GetV1CounterpartiesIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

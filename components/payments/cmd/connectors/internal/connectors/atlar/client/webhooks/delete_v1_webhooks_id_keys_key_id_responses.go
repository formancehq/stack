// Code generated by go-swagger; DO NOT EDIT.

package webhooks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/models"
)

// DeleteV1WebhooksIDKeysKeyIDReader is a Reader for the DeleteV1WebhooksIDKeysKeyID structure.
type DeleteV1WebhooksIDKeysKeyIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteV1WebhooksIDKeysKeyIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteV1WebhooksIDKeysKeyIDNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewDeleteV1WebhooksIDKeysKeyIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 412:
		result := NewDeleteV1WebhooksIDKeysKeyIDPreconditionFailed()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[DELETE /v1/webhooks/{id}/keys/{keyId}] DeleteV1WebhooksIDKeysKeyID", response, response.Code())
	}
}

// NewDeleteV1WebhooksIDKeysKeyIDNoContent creates a DeleteV1WebhooksIDKeysKeyIDNoContent with default headers values
func NewDeleteV1WebhooksIDKeysKeyIDNoContent() *DeleteV1WebhooksIDKeysKeyIDNoContent {
	return &DeleteV1WebhooksIDKeysKeyIDNoContent{}
}

/*
DeleteV1WebhooksIDKeysKeyIDNoContent describes a response with status code 204, with default header values.

No Content
*/
type DeleteV1WebhooksIDKeysKeyIDNoContent struct {
}

// IsSuccess returns true when this delete v1 webhooks Id keys key Id no content response has a 2xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this delete v1 webhooks Id keys key Id no content response has a 3xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete v1 webhooks Id keys key Id no content response has a 4xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) IsClientError() bool {
	return false
}

// IsServerError returns true when this delete v1 webhooks Id keys key Id no content response has a 5xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) IsServerError() bool {
	return false
}

// IsCode returns true when this delete v1 webhooks Id keys key Id no content response a status code equal to that given
func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) IsCode(code int) bool {
	return code == 204
}

// Code gets the status code for the delete v1 webhooks Id keys key Id no content response
func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) Code() int {
	return 204
}

func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) Error() string {
	return fmt.Sprintf("[DELETE /v1/webhooks/{id}/keys/{keyId}][%d] deleteV1WebhooksIdKeysKeyIdNoContent ", 204)
}

func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) String() string {
	return fmt.Sprintf("[DELETE /v1/webhooks/{id}/keys/{keyId}][%d] deleteV1WebhooksIdKeysKeyIdNoContent ", 204)
}

func (o *DeleteV1WebhooksIDKeysKeyIDNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteV1WebhooksIDKeysKeyIDNotFound creates a DeleteV1WebhooksIDKeysKeyIDNotFound with default headers values
func NewDeleteV1WebhooksIDKeysKeyIDNotFound() *DeleteV1WebhooksIDKeysKeyIDNotFound {
	return &DeleteV1WebhooksIDKeysKeyIDNotFound{}
}

/*
DeleteV1WebhooksIDKeysKeyIDNotFound describes a response with status code 404, with default header values.

The identified webhook key doesn't exist
*/
type DeleteV1WebhooksIDKeysKeyIDNotFound struct {
	Payload string
}

// IsSuccess returns true when this delete v1 webhooks Id keys key Id not found response has a 2xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete v1 webhooks Id keys key Id not found response has a 3xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete v1 webhooks Id keys key Id not found response has a 4xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete v1 webhooks Id keys key Id not found response has a 5xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this delete v1 webhooks Id keys key Id not found response a status code equal to that given
func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the delete v1 webhooks Id keys key Id not found response
func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) Code() int {
	return 404
}

func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) Error() string {
	return fmt.Sprintf("[DELETE /v1/webhooks/{id}/keys/{keyId}][%d] deleteV1WebhooksIdKeysKeyIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) String() string {
	return fmt.Sprintf("[DELETE /v1/webhooks/{id}/keys/{keyId}][%d] deleteV1WebhooksIdKeysKeyIdNotFound  %+v", 404, o.Payload)
}

func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) GetPayload() string {
	return o.Payload
}

func (o *DeleteV1WebhooksIDKeysKeyIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteV1WebhooksIDKeysKeyIDPreconditionFailed creates a DeleteV1WebhooksIDKeysKeyIDPreconditionFailed with default headers values
func NewDeleteV1WebhooksIDKeysKeyIDPreconditionFailed() *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed {
	return &DeleteV1WebhooksIDKeysKeyIDPreconditionFailed{}
}

/*
DeleteV1WebhooksIDKeysKeyIDPreconditionFailed describes a response with status code 412, with default header values.

ErrorResponse with code `key_required`
*/
type DeleteV1WebhooksIDKeysKeyIDPreconditionFailed struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this delete v1 webhooks Id keys key Id precondition failed response has a 2xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this delete v1 webhooks Id keys key Id precondition failed response has a 3xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) IsRedirect() bool {
	return false
}

// IsClientError returns true when this delete v1 webhooks Id keys key Id precondition failed response has a 4xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) IsClientError() bool {
	return true
}

// IsServerError returns true when this delete v1 webhooks Id keys key Id precondition failed response has a 5xx status code
func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) IsServerError() bool {
	return false
}

// IsCode returns true when this delete v1 webhooks Id keys key Id precondition failed response a status code equal to that given
func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) IsCode(code int) bool {
	return code == 412
}

// Code gets the status code for the delete v1 webhooks Id keys key Id precondition failed response
func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) Code() int {
	return 412
}

func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) Error() string {
	return fmt.Sprintf("[DELETE /v1/webhooks/{id}/keys/{keyId}][%d] deleteV1WebhooksIdKeysKeyIdPreconditionFailed  %+v", 412, o.Payload)
}

func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) String() string {
	return fmt.Sprintf("[DELETE /v1/webhooks/{id}/keys/{keyId}][%d] deleteV1WebhooksIdKeysKeyIdPreconditionFailed  %+v", 412, o.Payload)
}

func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DeleteV1WebhooksIDKeysKeyIDPreconditionFailed) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

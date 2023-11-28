// Code generated by go-swagger; DO NOT EDIT.

package expected_transactions

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/models"
)

// GetV1ExpectedTransactionsReader is a Reader for the GetV1ExpectedTransactions structure.
type GetV1ExpectedTransactionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetV1ExpectedTransactionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetV1ExpectedTransactionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("[GET /v1/expected-transactions] GetV1ExpectedTransactions", response, response.Code())
	}
}

// NewGetV1ExpectedTransactionsOK creates a GetV1ExpectedTransactionsOK with default headers values
func NewGetV1ExpectedTransactionsOK() *GetV1ExpectedTransactionsOK {
	return &GetV1ExpectedTransactionsOK{}
}

/*
GetV1ExpectedTransactionsOK describes a response with status code 200, with default header values.

QueryResponse with list of expected transactions
*/
type GetV1ExpectedTransactionsOK struct {
	Payload *GetV1ExpectedTransactionsOKBody
}

// IsSuccess returns true when this get v1 expected transactions o k response has a 2xx status code
func (o *GetV1ExpectedTransactionsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get v1 expected transactions o k response has a 3xx status code
func (o *GetV1ExpectedTransactionsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 expected transactions o k response has a 4xx status code
func (o *GetV1ExpectedTransactionsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get v1 expected transactions o k response has a 5xx status code
func (o *GetV1ExpectedTransactionsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get v1 expected transactions o k response a status code equal to that given
func (o *GetV1ExpectedTransactionsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get v1 expected transactions o k response
func (o *GetV1ExpectedTransactionsOK) Code() int {
	return 200
}

func (o *GetV1ExpectedTransactionsOK) Error() string {
	return fmt.Sprintf("[GET /v1/expected-transactions][%d] getV1ExpectedTransactionsOK  %+v", 200, o.Payload)
}

func (o *GetV1ExpectedTransactionsOK) String() string {
	return fmt.Sprintf("[GET /v1/expected-transactions][%d] getV1ExpectedTransactionsOK  %+v", 200, o.Payload)
}

func (o *GetV1ExpectedTransactionsOK) GetPayload() *GetV1ExpectedTransactionsOKBody {
	return o.Payload
}

func (o *GetV1ExpectedTransactionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetV1ExpectedTransactionsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetV1ExpectedTransactionsOKBody get v1 expected transactions o k body
swagger:model GetV1ExpectedTransactionsOKBody
*/
type GetV1ExpectedTransactionsOKBody struct {
	models.QueryResponse

	// items
	Items []*models.ExpectedTransaction `json:"items"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *GetV1ExpectedTransactionsOKBody) UnmarshalJSON(raw []byte) error {
	// GetV1ExpectedTransactionsOKBodyAO0
	var getV1ExpectedTransactionsOKBodyAO0 models.QueryResponse
	if err := swag.ReadJSON(raw, &getV1ExpectedTransactionsOKBodyAO0); err != nil {
		return err
	}
	o.QueryResponse = getV1ExpectedTransactionsOKBodyAO0

	// GetV1ExpectedTransactionsOKBodyAO1
	var dataGetV1ExpectedTransactionsOKBodyAO1 struct {
		Items []*models.ExpectedTransaction `json:"items"`
	}
	if err := swag.ReadJSON(raw, &dataGetV1ExpectedTransactionsOKBodyAO1); err != nil {
		return err
	}

	o.Items = dataGetV1ExpectedTransactionsOKBodyAO1.Items

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o GetV1ExpectedTransactionsOKBody) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	getV1ExpectedTransactionsOKBodyAO0, err := swag.WriteJSON(o.QueryResponse)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, getV1ExpectedTransactionsOKBodyAO0)
	var dataGetV1ExpectedTransactionsOKBodyAO1 struct {
		Items []*models.ExpectedTransaction `json:"items"`
	}

	dataGetV1ExpectedTransactionsOKBodyAO1.Items = o.Items

	jsonDataGetV1ExpectedTransactionsOKBodyAO1, errGetV1ExpectedTransactionsOKBodyAO1 := swag.WriteJSON(dataGetV1ExpectedTransactionsOKBodyAO1)
	if errGetV1ExpectedTransactionsOKBodyAO1 != nil {
		return nil, errGetV1ExpectedTransactionsOKBodyAO1
	}
	_parts = append(_parts, jsonDataGetV1ExpectedTransactionsOKBodyAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this get v1 expected transactions o k body
func (o *GetV1ExpectedTransactionsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with models.QueryResponse
	if err := o.QueryResponse.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateItems(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetV1ExpectedTransactionsOKBody) validateItems(formats strfmt.Registry) error {

	if swag.IsZero(o.Items) { // not required
		return nil
	}

	for i := 0; i < len(o.Items); i++ {
		if swag.IsZero(o.Items[i]) { // not required
			continue
		}

		if o.Items[i] != nil {
			if err := o.Items[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getV1ExpectedTransactionsOK" + "." + "items" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getV1ExpectedTransactionsOK" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get v1 expected transactions o k body based on the context it is used
func (o *GetV1ExpectedTransactionsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with models.QueryResponse
	if err := o.QueryResponse.ContextValidate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := o.contextValidateItems(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetV1ExpectedTransactionsOKBody) contextValidateItems(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Items); i++ {

		if o.Items[i] != nil {

			if swag.IsZero(o.Items[i]) { // not required
				return nil
			}

			if err := o.Items[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getV1ExpectedTransactionsOK" + "." + "items" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getV1ExpectedTransactionsOK" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetV1ExpectedTransactionsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetV1ExpectedTransactionsOKBody) UnmarshalBinary(b []byte) error {
	var res GetV1ExpectedTransactionsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
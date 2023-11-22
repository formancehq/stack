// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TestbankReturnReason testbank return reason
//
// swagger:model testbank.ReturnReason
type TestbankReturnReason struct {

	// code
	Code string `json:"code,omitempty"`
}

// Validate validates this testbank return reason
func (m *TestbankReturnReason) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this testbank return reason based on context it is used
func (m *TestbankReturnReason) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TestbankReturnReason) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TestbankReturnReason) UnmarshalBinary(b []byte) error {
	var res TestbankReturnReason
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

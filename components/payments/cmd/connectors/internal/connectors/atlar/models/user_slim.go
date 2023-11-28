// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UserSlim user slim
//
// swagger:model UserSlim
type UserSlim struct {

	// first name
	// Example: John
	FirstName string `json:"firstName,omitempty"`

	// id
	// Example: 54528335-739e-4de9-bdb3-28f96c98785d
	ID string `json:"id,omitempty"`

	// last name
	// Example: Smith
	LastName string `json:"lastName,omitempty"`

	// programmatic access
	// Example: false
	ProgrammaticAccess bool `json:"programmaticAccess,omitempty"`

	// username
	// Example: john.smith@johnnysmithson.com
	Username string `json:"username,omitempty"`
}

// Validate validates this user slim
func (m *UserSlim) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this user slim based on context it is used
func (m *UserSlim) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UserSlim) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UserSlim) UnmarshalBinary(b []byte) error {
	var res UserSlim
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
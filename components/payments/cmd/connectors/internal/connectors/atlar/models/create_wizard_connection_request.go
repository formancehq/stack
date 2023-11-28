// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CreateWizardConnectionRequest create wizard connection request
//
// swagger:model CreateWizardConnectionRequest
type CreateWizardConnectionRequest struct {

	// connection type Id
	ConnectionTypeID string `json:"connectionTypeId,omitempty"`
}

// Validate validates this create wizard connection request
func (m *CreateWizardConnectionRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this create wizard connection request based on context it is used
func (m *CreateWizardConnectionRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreateWizardConnectionRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateWizardConnectionRequest) UnmarshalBinary(b []byte) error {
	var res CreateWizardConnectionRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
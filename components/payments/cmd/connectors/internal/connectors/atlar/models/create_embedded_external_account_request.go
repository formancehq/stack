// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CreateEmbeddedExternalAccountRequest create embedded external account request
//
// swagger:model CreateEmbeddedExternalAccountRequest
type CreateEmbeddedExternalAccountRequest struct {

	// bank
	Bank *UpdatableBank `json:"bank,omitempty"`

	// Deprecated. Use bank.bic
	Bic string `json:"bic,omitempty"`

	// ExternalId is optional to use, but if used, the Atlar platform will persist it, index it, as well as require it to be unique. It is also possible to retrieve the identified resource using the ExternalId.
	// Example: walVNuin6X5Mvte4xhg1ibTAVSACfN4Q9hl
	ExternalID string `json:"externalId,omitempty"`

	// Any external metadata you want to attach, such as your own internal IDs.
	ExternalMetadata ExternalMetadata `json:"externalMetadata,omitempty"`

	// A list of account identifiers with account numbers that identify this account. A single account identifier can be enough, but some countries use multiple. For example, an account can be described both with an IBAN and one or more local formats.
	// Required: true
	Identifiers []*AccountIdentifier `json:"identifiers"`
}

// Validate validates this create embedded external account request
func (m *CreateEmbeddedExternalAccountRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBank(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExternalMetadata(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIdentifiers(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreateEmbeddedExternalAccountRequest) validateBank(formats strfmt.Registry) error {
	if swag.IsZero(m.Bank) { // not required
		return nil
	}

	if m.Bank != nil {
		if err := m.Bank.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bank")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("bank")
			}
			return err
		}
	}

	return nil
}

func (m *CreateEmbeddedExternalAccountRequest) validateExternalMetadata(formats strfmt.Registry) error {
	if swag.IsZero(m.ExternalMetadata) { // not required
		return nil
	}

	if m.ExternalMetadata != nil {
		if err := m.ExternalMetadata.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("externalMetadata")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("externalMetadata")
			}
			return err
		}
	}

	return nil
}

func (m *CreateEmbeddedExternalAccountRequest) validateIdentifiers(formats strfmt.Registry) error {

	if err := validate.Required("identifiers", "body", m.Identifiers); err != nil {
		return err
	}

	for i := 0; i < len(m.Identifiers); i++ {
		if swag.IsZero(m.Identifiers[i]) { // not required
			continue
		}

		if m.Identifiers[i] != nil {
			if err := m.Identifiers[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("identifiers" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("identifiers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this create embedded external account request based on the context it is used
func (m *CreateEmbeddedExternalAccountRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBank(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateExternalMetadata(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIdentifiers(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreateEmbeddedExternalAccountRequest) contextValidateBank(ctx context.Context, formats strfmt.Registry) error {

	if m.Bank != nil {

		if swag.IsZero(m.Bank) { // not required
			return nil
		}

		if err := m.Bank.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bank")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("bank")
			}
			return err
		}
	}

	return nil
}

func (m *CreateEmbeddedExternalAccountRequest) contextValidateExternalMetadata(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.ExternalMetadata) { // not required
		return nil
	}

	if err := m.ExternalMetadata.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("externalMetadata")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("externalMetadata")
		}
		return err
	}

	return nil
}

func (m *CreateEmbeddedExternalAccountRequest) contextValidateIdentifiers(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Identifiers); i++ {

		if m.Identifiers[i] != nil {

			if swag.IsZero(m.Identifiers[i]) { // not required
				return nil
			}

			if err := m.Identifiers[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("identifiers" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("identifiers" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *CreateEmbeddedExternalAccountRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateEmbeddedExternalAccountRequest) UnmarshalBinary(b []byte) error {
	var res CreateEmbeddedExternalAccountRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

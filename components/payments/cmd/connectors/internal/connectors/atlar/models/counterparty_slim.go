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

// CounterpartySlim counterparty slim
//
// swagger:model CounterpartySlim
type CounterpartySlim struct {

	// contact details
	ContactDetails *ContactDetailsSlim `json:"contactDetails,omitempty"`

	// id
	// Example: f3efbb73-4e5b-4b22-adeb-918bbf1dfbd8
	// Required: true
	ID *string `json:"id"`

	// Deprecated. Has moved to external accounts. A list of AccountIdentifiers. An AccountIdentifier uniquely identifies one bank account.
	Identifiers []*AccountIdentifierSlim `json:"identifiers"`

	// Name of your Counterparty
	// Example: Customer #312
	Name string `json:"name,omitempty"`

	// The legal type of the Counterparty
	// Example: INDIVIDUAL
	PartyType string `json:"partyType,omitempty"`
}

// Validate validates this counterparty slim
func (m *CounterpartySlim) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContactDetails(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
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

func (m *CounterpartySlim) validateContactDetails(formats strfmt.Registry) error {
	if swag.IsZero(m.ContactDetails) { // not required
		return nil
	}

	if m.ContactDetails != nil {
		if err := m.ContactDetails.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("contactDetails")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("contactDetails")
			}
			return err
		}
	}

	return nil
}

func (m *CounterpartySlim) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *CounterpartySlim) validateIdentifiers(formats strfmt.Registry) error {
	if swag.IsZero(m.Identifiers) { // not required
		return nil
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

// ContextValidate validate this counterparty slim based on the context it is used
func (m *CounterpartySlim) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateContactDetails(ctx, formats); err != nil {
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

func (m *CounterpartySlim) contextValidateContactDetails(ctx context.Context, formats strfmt.Registry) error {

	if m.ContactDetails != nil {

		if swag.IsZero(m.ContactDetails) { // not required
			return nil
		}

		if err := m.ContactDetails.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("contactDetails")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("contactDetails")
			}
			return err
		}
	}

	return nil
}

func (m *CounterpartySlim) contextValidateIdentifiers(ctx context.Context, formats strfmt.Registry) error {

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
func (m *CounterpartySlim) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CounterpartySlim) UnmarshalBinary(b []byte) error {
	var res CounterpartySlim
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

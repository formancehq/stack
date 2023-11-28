// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TriggerCondition trigger condition
//
// swagger:model TriggerCondition
type TriggerCondition struct {

	// Payments with an amount value less than this value will trigger the chain.
	AmountLt *Amount `json:"amountLt,omitempty"`

	// Payments created by this (optional) role ID will trigger the chain.
	CreatorRoleID string `json:"creatorRoleId,omitempty"`
}

// Validate validates this trigger condition
func (m *TriggerCondition) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAmountLt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TriggerCondition) validateAmountLt(formats strfmt.Registry) error {
	if swag.IsZero(m.AmountLt) { // not required
		return nil
	}

	if m.AmountLt != nil {
		if err := m.AmountLt.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("amountLt")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("amountLt")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this trigger condition based on the context it is used
func (m *TriggerCondition) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAmountLt(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TriggerCondition) contextValidateAmountLt(ctx context.Context, formats strfmt.Registry) error {

	if m.AmountLt != nil {

		if swag.IsZero(m.AmountLt) { // not required
			return nil
		}

		if err := m.AmountLt.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("amountLt")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("amountLt")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TriggerCondition) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TriggerCondition) UnmarshalBinary(b []byte) error {
	var res TriggerCondition
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
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

// TransactionCharacteristics transaction characteristics
//
// swagger:model TransactionCharacteristics
type TransactionCharacteristics struct {

	// ISO 20022 Bank Transaction Code.
	BankTransactionCode *BankTransactionCode `json:"bankTransactionCode,omitempty"`

	// Optional. If any currency conversion happened at the bank holding the account, conversion details may be present here.
	CurrencyExchange *CurrencyExchange `json:"currencyExchange,omitempty"`

	// The amount that the initiator of the transaction instructed.
	InstructedAmount *Amount `json:"instructedAmount,omitempty"`

	// ISO 20022 Return Reason.
	ReturnReason *ReturnReason `json:"returnReason,omitempty"`

	// Whether or not the transaction is considered a return or reversal.
	Returned bool `json:"returned,omitempty"`
}

// Validate validates this transaction characteristics
func (m *TransactionCharacteristics) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBankTransactionCode(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCurrencyExchange(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateInstructedAmount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateReturnReason(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TransactionCharacteristics) validateBankTransactionCode(formats strfmt.Registry) error {
	if swag.IsZero(m.BankTransactionCode) { // not required
		return nil
	}

	if m.BankTransactionCode != nil {
		if err := m.BankTransactionCode.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bankTransactionCode")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("bankTransactionCode")
			}
			return err
		}
	}

	return nil
}

func (m *TransactionCharacteristics) validateCurrencyExchange(formats strfmt.Registry) error {
	if swag.IsZero(m.CurrencyExchange) { // not required
		return nil
	}

	if m.CurrencyExchange != nil {
		if err := m.CurrencyExchange.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("currencyExchange")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("currencyExchange")
			}
			return err
		}
	}

	return nil
}

func (m *TransactionCharacteristics) validateInstructedAmount(formats strfmt.Registry) error {
	if swag.IsZero(m.InstructedAmount) { // not required
		return nil
	}

	if m.InstructedAmount != nil {
		if err := m.InstructedAmount.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("instructedAmount")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("instructedAmount")
			}
			return err
		}
	}

	return nil
}

func (m *TransactionCharacteristics) validateReturnReason(formats strfmt.Registry) error {
	if swag.IsZero(m.ReturnReason) { // not required
		return nil
	}

	if m.ReturnReason != nil {
		if err := m.ReturnReason.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("returnReason")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("returnReason")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this transaction characteristics based on the context it is used
func (m *TransactionCharacteristics) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBankTransactionCode(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateCurrencyExchange(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateInstructedAmount(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateReturnReason(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TransactionCharacteristics) contextValidateBankTransactionCode(ctx context.Context, formats strfmt.Registry) error {

	if m.BankTransactionCode != nil {

		if swag.IsZero(m.BankTransactionCode) { // not required
			return nil
		}

		if err := m.BankTransactionCode.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bankTransactionCode")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("bankTransactionCode")
			}
			return err
		}
	}

	return nil
}

func (m *TransactionCharacteristics) contextValidateCurrencyExchange(ctx context.Context, formats strfmt.Registry) error {

	if m.CurrencyExchange != nil {

		if swag.IsZero(m.CurrencyExchange) { // not required
			return nil
		}

		if err := m.CurrencyExchange.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("currencyExchange")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("currencyExchange")
			}
			return err
		}
	}

	return nil
}

func (m *TransactionCharacteristics) contextValidateInstructedAmount(ctx context.Context, formats strfmt.Registry) error {

	if m.InstructedAmount != nil {

		if swag.IsZero(m.InstructedAmount) { // not required
			return nil
		}

		if err := m.InstructedAmount.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("instructedAmount")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("instructedAmount")
			}
			return err
		}
	}

	return nil
}

func (m *TransactionCharacteristics) contextValidateReturnReason(ctx context.Context, formats strfmt.Registry) error {

	if m.ReturnReason != nil {

		if swag.IsZero(m.ReturnReason) { // not required
			return nil
		}

		if err := m.ReturnReason.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("returnReason")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("returnReason")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TransactionCharacteristics) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TransactionCharacteristics) UnmarshalBinary(b []byte) error {
	var res TransactionCharacteristics
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

package wallet

import (
	"errors"
	"fmt"
)

var (
	ErrWalletNotFound        = errors.New("wallet_not_found")
	ErrHoldNotFound          = errors.New("hold_not_found")
	ErrInsufficientFundError = errors.New("insufficient fund")
	ErrClosedHold            = errors.New("closed hold")
)

type MismatchTypeError struct {
	expected, got string
}

func (t MismatchTypeError) Error() string {
	return fmt.Sprintf("unexpected type, got '%s', but '%s' was expected", t.got, t.expected)
}

func newErrMismatchType(expected, got string) MismatchTypeError {
	return MismatchTypeError{
		expected: expected,
		got:      got,
	}
}

func IsMismatchTypeError(err error) bool {
	return errors.As(err, &MismatchTypeError{})
}

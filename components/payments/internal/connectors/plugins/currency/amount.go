package currency

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

// This package provides a way to convert a string amount to a big.Int
// with a given precision. The precision is the number of decimals that
// the amount has. For example, if the amount is 123.45 and the precision is 2,
// the amount will be 12345.
// We also provide a way to convert a big.Int to a string amount with a given
// precision.
// We developed this package because we need to convert amounts from strings,
// and apply the precision to them to have minor units. If we do that with
// the big package using floats and divisions/multiplcations, we will
// sometimes loose precision, which is unacceptable.
// Here is an example of loosing precision with the big package:

var (
	// ErrInvalidAmount is returned when the amount is invalid
	ErrInvalidAmount = fmt.Errorf("invalid amount")
	// ErrInvalidPrecision is returned when the precision is inferior to the
	// number of decimals in the amount or negative
	ErrInvalidPrecision = fmt.Errorf("invalid precision")
)

func GetAmountWithPrecisionFromString(amountString string, precision int) (*big.Int, error) {
	if precision < 0 {
		return nil, errors.Wrap(ErrInvalidPrecision, fmt.Sprintf("precision is negative: %d", precision))
	}

	parts := strings.Split(amountString, ".")

	lengthParts := len(parts)

	if lengthParts > 2 || lengthParts == 0 {
		// More than one dot, invalid amount
		return nil, errors.Wrap(ErrInvalidAmount, fmt.Sprintf("got multiple dots in amount: %s", amountString))
	}

	if lengthParts == 1 {
		// No dot, which means it's an integer
		for p := 0; p < precision; p++ {
			amountString += "0"
		}
		res, ok := new(big.Int).SetString(amountString, 10)
		if !ok {
			return nil, errors.Wrap(ErrInvalidAmount, fmt.Sprintf("invalid amount: %s", amountString))
		}
		return res, nil
	}

	// Here we are in the case where we have one dot, which means we have a
	// decimal amount
	decimalPart := parts[1]
	lengthDecimalPart := len(decimalPart)
	switch {
	case lengthDecimalPart == precision:
		// The decimal part has the same length as the precision, we can
		// concatenate the two parts and return the result
		res, ok := new(big.Int).SetString(parts[0]+decimalPart, 10)
		if !ok {
			return nil, errors.Wrap(ErrInvalidAmount, fmt.Sprintf("invalid amount computed: %s from amount %s", parts[0]+decimalPart, amountString))
		}
		return res, nil

	case lengthDecimalPart < precision:
		// The decimal part is shorter than the precision, we need to add
		// some zeros at the end of the decimal part
		for p := 0; p < precision-lengthDecimalPart; p++ {
			decimalPart += "0"
		}
		res, ok := new(big.Int).SetString(parts[0]+decimalPart, 10)
		if !ok {
			return nil, errors.Wrap(ErrInvalidAmount, fmt.Sprintf("invalid amount computed: %s from amount %s", parts[0]+decimalPart, amountString))
		}
		return res, nil

	default:
		// The decimal part is longer than the precision, we have to send an
		// error because we don't want to loose the precision
		return nil, ErrInvalidPrecision
	}
}

func GetStringAmountFromBigIntWithPrecision(amount *big.Int, precision int) (string, error) {
	if precision < 0 {
		return "", errors.Wrap(ErrInvalidPrecision, fmt.Sprintf("precision is negative: %d", precision))
	}

	amountString := amount.String()
	amountStringLength := len(amountString)

	if precision == 0 {
		// Nothing to do
		return amountString, nil
	}

	decimalPart := bytes.NewBufferString("")
	for p := precision; p > 0; p-- {
		if amountStringLength < p {
			decimalPart.WriteByte('0')
			continue
		}
		decimalPart.WriteByte(amountString[amountStringLength-p])
	}

	if amountStringLength < precision || amountStringLength == precision {
		return "0." + decimalPart.String(), nil
	}

	// Here we are in the case where the amount has more digits than the
	// precision, we need to add a dot at the right place
	return amountString[:amountStringLength-precision] + "." + decimalPart.String(), nil
}

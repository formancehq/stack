package atlar

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtlarTransactionAmountToPaymentAbsoluteAmount(t *testing.T) {
	t.Parallel()
	var result *big.Int
	var err error

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount("0.30", 2)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(30), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount("3.30", 2)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(330), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount("33.0", 1)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(330), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount("-0.30", 2)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(30), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount("-3.30", 2)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(330), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount("hello world", 2)
	if assert.Error(t, err, "failed to parse amount hello world") {
		assert.Nil(t, result)
	}
}

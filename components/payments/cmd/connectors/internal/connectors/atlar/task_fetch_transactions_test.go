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

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount(30)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(30), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount(330)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(330), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount(330)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(330), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount(-30)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(30), *result)
	}

	result, err = atlarTransactionAmountToPaymentAbsoluteAmount(-330)
	if assert.Nil(t, err) {
		assert.Equal(t, *big.NewInt(330), *result)
	}
}

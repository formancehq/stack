package atlar

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAmountToString(t *testing.T) {
	t.Parallel()

	assert.EqualValues(t, "0.032", amountToString(*big.NewInt(32), 3))
	assert.EqualValues(t, "0.32", amountToString(*big.NewInt(32), 2))
	assert.EqualValues(t, "3.2", amountToString(*big.NewInt(32), 1))
	assert.EqualValues(t, "5.432", amountToString(*big.NewInt(5432), 3))
	assert.EqualValues(t, "54.32", amountToString(*big.NewInt(5432), 2))
	assert.EqualValues(t, "543.2", amountToString(*big.NewInt(5432), 1))
}

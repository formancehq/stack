package atlar

import (
	"errors"
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

func TestSerializeAtlarPaymentExternalID(t *testing.T) {
	t.Parallel()

	assert.EqualValues(t, "testID_1", serializeAtlarPaymentExternalID("testID", 1))
	assert.EqualValues(t, "tqmbAGgV4S2pHics57BT5tV2_682", serializeAtlarPaymentExternalID("tqmbAGgV4S2pHics57BT5tV2", 682))
}

func TestDeserializeAtlarPaymentExternalID(t *testing.T) {
	t.Parallel()

	var ID string
	var attempts int
	var err error

	ID, attempts, err = deserializeAtlarPaymentExternalID("testID_1")
	if assert.Nil(t, err) {
		assert.EqualValues(t, "testID", ID)
		assert.EqualValues(t, 1, attempts)
	}

	ID, attempts, err = deserializeAtlarPaymentExternalID("tqmbAGgV4S2pHics57BT5tV2_682")
	if assert.Nil(t, err) {
		assert.EqualValues(t, "tqmbAGgV4S2pHics57BT5tV2", ID)
		assert.EqualValues(t, 682, attempts)
	}

	_, _, err = deserializeAtlarPaymentExternalID("tqmbAGgV4S2pHics57BT5tV2_682_432")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("cannot deserialize malformed externalID"), err)
	}

	_, _, err = deserializeAtlarPaymentExternalID("tqmbAGgV4S2pHics57BT5tV2_")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("cannot deserialize malformed externalID"), err)
	}

	_, _, err = deserializeAtlarPaymentExternalID("tqmbAGgV4S2pHics57BT5tV2")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("cannot deserialize malformed externalID"), err)
	}
}

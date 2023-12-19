package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataToIdentifierData(t *testing.T) {
	t.Parallel()

	_, err := metadataToIdentifierData("not_valid", "test")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("input does not match the expected format"), err)
	}
	_, err = metadataToIdentifierData(atlarMetadataSpecNamespace+"not_valid", "test")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("input does not match the expected format"), err)
	}
	_, err = metadataToIdentifierData(atlarMetadataSpecNamespace+"identifier/not_valid", "test")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("input does not match the expected format"), err)
	}
	identifier, err := metadataToIdentifierData(atlarMetadataSpecNamespace+"identifier/DE/IBAN", "DE02700100800030876808")
	if assert.Nil(t, err) {
		assert.Equal(t, IdentifierData{
			Market: "DE",
			Type:   "IBAN",
			Number: "DE02700100800030876808",
		}, *identifier)
	}
}

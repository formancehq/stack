package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	assert.NoError(t, Config{
		Endpoint:   "https://example.com",
		EventTypes: []string{"TYPE1", "TYPE2"},
	}.Validate())

	assert.Error(t, Config{
		Endpoint:   " http://invalid",
		EventTypes: []string{"TYPE1", "TYPE2"},
	}.Validate())

	assert.Error(t, Config{
		Endpoint:   "https://example.com",
		EventTypes: []string{"TYPE1", ""},
	}.Validate())
}

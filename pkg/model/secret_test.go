package model

import (
	"encoding/base64"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecret_Validate(t *testing.T) {
	sec := Secret{Secret: NewSecret()}
	assert.NoError(t, sec.Validate())

	sec = Secret{}
	assert.NoError(t, sec.Validate())

	sec = Secret{Secret: "invalid"}
	assert.Error(t, sec.Validate())

	sec = Secret{Secret: base64.StdEncoding.EncodeToString([]byte(`invalid`))}
	assert.Error(t, sec.Validate())

	token := make([]byte, 23)
	rand.Read(token)
	tooShort := base64.StdEncoding.EncodeToString(token)
	sec = Secret{Secret: tooShort}
	assert.Error(t, sec.Validate())

	token = make([]byte, 25)
	rand.Read(token)
	tooLong := base64.StdEncoding.EncodeToString(token)
	sec = Secret{Secret: tooLong}
	assert.Error(t, sec.Validate())
}

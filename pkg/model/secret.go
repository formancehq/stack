package model

import (
	"encoding/base64"
	"fmt"
	"math/rand"
)

type Secret struct {
	Secret string `json:"secret" bson:"secret"`
}

func (s *Secret) Validate() error {
	if s.Secret == "" {
		s.Secret = NewSecret()
	} else {
		var decoded []byte
		var err error
		if decoded, err = base64.StdEncoding.DecodeString(s.Secret); err != nil {
			return fmt.Errorf("secret should be base64 encoded: %w", err)
		}
		if len(decoded) != 24 {
			return fmt.Errorf("decoded secret should be of size 24: actual size %d", len(decoded))
		}
	}

	return nil
}

func NewSecret() string {
	token := make([]byte, 24)
	rand.Read(token)
	return base64.StdEncoding.EncodeToString(token)
}

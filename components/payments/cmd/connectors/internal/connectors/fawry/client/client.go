package client

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type Client struct {
	SecureKey    string
	MerchantCode string
}

func (c *Client) SignRequest(inputs []string) string {
	inputs = append(inputs, c.SecureKey)
	sum := sha256.Sum256([]byte(strings.Join(inputs[:], ",")))
	return hex.EncodeToString(sum[:])
}

func (c *Client) Init() error {
	return nil
}

func NewClient() *Client {
	return &Client{}
}

package client

import (
	"encoding/json"
	"fmt"
	"io"
)

type mangopayError struct {
	StatusCode int               `json:"-"`
	Message    string            `json:"Message"`
	Type       string            `json:"Type"`
	Errors     map[string]string `json:"Errors"`
}

func (me *mangopayError) Error() error {
	var errorMessage string
	if len(me.Errors) > 0 {
		for _, message := range me.Errors {
			errorMessage = message
			break
		}
	}

	if errorMessage == "" {
		return fmt.Errorf("unexpected status code: %d", me.StatusCode)
	}

	return fmt.Errorf("%d: %s", me.StatusCode, errorMessage)
}

func unmarshalError(statusCode int, body io.ReadCloser) *mangopayError {
	var ce mangopayError
	_ = json.NewDecoder(body).Decode(&ce)

	ce.StatusCode = statusCode

	return &ce
}
